package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/jinzhu/gorm"
)

type Task struct {
	Id       string
	db       *gorm.DB
	guard    *Guard
	success  int
	fail     int
	target   string
	stop     chan int // 0 stop, 1 cancel
	quitChan chan string
}

func NewTask(quitChan chan string, db *gorm.DB) (*Task, error) {
	fileModel := entity.FileModel{}
	if err := db.Where("status = ?", entity.Waitting).First(&fileModel).Error; err != nil {
		return nil, err
	}
	batchMode := entity.BatchModel{}
	if err := db.Where("id = ?", fileModel.BatchId).First(&batchMode).Error; err != nil {
		return nil, err
	}
	db.Model(&fileModel).Update("status", entity.Processing)
	return &Task{
			Id:       fileModel.Id,
			db:       db,
			target:   batchMode.Api,
			quitChan: quitChan,
			stop:     make(chan int, 1),
		},
		nil
}

func (task *Task) Run() error {
	g, s, f, err := GuardOpen(task.Id, task.db)
	if err != nil {
		// TODO
		return err
	}
	task.guard = g
	task.success = s
	task.fail = f
	go task.run()
	return nil
}

func (task *Task) Stop() {
	task.stop <- 0
}

func (task *Task) Cancel() {
	task.stop <- 1
}

func (task *Task) Pause() {
	task.stop <- 2
}

func (task *Task) Info() (info sth.Result, err error) {
	if info, err = task.guard.Info(); err == nil {
		info["success"] = task.success
		info["fail"] = task.fail
	}
	return
}

func (task *Task) run() {
	var stopAt = entity.Processing
loop:
	for i := 0; ; i++ {
		select {
		case code := <-task.stop:
			if code == 0 {
				stopAt = entity.Waitting
			} else if code == 1 {
				stopAt = entity.Cancel
			} else {
				stopAt = entity.Paused
			}
			break loop
		default:
			if i >= 10 {
				i = 0
				task.update()
			}
			if err := task.sendOne(); err != nil {
				if err == io.EOF {
					stopAt = entity.Finish
				} else if err == SendError {
					stopAt = entity.Fail
				}
				break loop
			}
		}
	}
	task.update()
	task.guard.GuardClose(stopAt)
	task.quitChan <- task.Id
}

func (task *Task) update() {
	task.guard.Update(task.success, task.fail)
}

var SendError = errors.New("Server unavailable")

func (task *Task) sendOne() error {
	line, err := task.guard.ReadLine()
	line = strings.TrimSpace(line)
	if len(line) > 0 {
		jsonValue, _ := json.Marshal(map[string]interface{}{"url": line})
		for i := 0; i < 3; i++ {
			resp, err := http.Post(task.target, "application/json", bytes.NewBuffer(jsonValue))
			if err == nil {
				logger.Debug("Send", line, resp.StatusCode)
				if resp.StatusCode == http.StatusOK {
					task.success += 1
				} else if resp.StatusCode == http.StatusBadRequest {
					task.fail += 1
				}
				return err
			} else {
				time.Sleep(time.Second)
				continue
			}
		}
		return SendError
	}
	return err
}
