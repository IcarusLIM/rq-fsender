package task

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type Task struct {
	Id         string
	db         *gorm.DB
	lock       *sync.RWMutex
	fm         *entity.FileModel
	guard      *Guard
	eofFlag    bool
	success    int32
	fail       int32
	concur     int32
	reciver    string
	logPath    string
	stop       chan int32 // signal from taskbox, trigger by user: 0 stop, 1 cancel, 2 pasue
	reportQuit chan string
}

func NewTask(quitChan chan string, db *gorm.DB, conf *viper.Viper) (*Task, error) {
	fileModel := entity.FileModel{}
	// TODO transaction control
	if err := db.Where("status = ? OR (status = ? AND update_at < ?)",
		entity.Waitting,
		entity.Processing,
		time.Now().Add(-time.Minute*2)).First(&fileModel).Error; err != nil {
		return nil, err
	}
	batchMode := entity.BatchModel{}
	if err := db.Where("id = ?", fileModel.BatchId).First(&batchMode).Error; err != nil {
		return nil, err
	}
	db.Model(&fileModel).Update("status", entity.Processing)
	concur := conf.GetInt32("task.conf.concurrent")

	return &Task{
			Id:         fileModel.Id,
			db:         db,
			lock:       &sync.RWMutex{},
			fm:         &fileModel,
			eofFlag:    false,
			success:    fileModel.Success,
			fail:       fileModel.Fail,
			concur:     concur,
			reciver:    batchMode.Api,
			logPath:    conf.GetString("upload.log"),
			stop:       make(chan int32, concur+1),
			reportQuit: quitChan,
		},
		nil
}

func (task *Task) Run() error {
	var fileMeta entity.FileMeta
	if err := json.Unmarshal([]byte(task.fm.FileMeta), &fileMeta); err != nil {
		return err
	}
	g, err := GuardOpen(fileMeta, task.fm.Offset, task.db)
	if err != nil {
		// TODO
		return err
	}
	task.guard = g
	go task.run()
	return nil
}

func (task *Task) run() {
	concur := 10
	respChan := make(chan map[string]interface{}, concur)
	closeWait := make(chan struct{})
	loopWG := sync.WaitGroup{}
	loopWG.Add(concur)
	for i := 0; i < concur; i++ {
		go task.sendloop(respChan, closeWait, &loopWG)
	}
	loopWGChan := make(chan struct{})
	go func() {
		loopWG.Wait()
		close(loopWGChan)
	}()

	postWG := sync.WaitGroup{}
	postWG.Add(1)
	go task.postSend(&postWG, respChan)

	stopAt := entity.Processing
	select {
	case code := <-task.stop:
		if code == 0 {
			stopAt = entity.Waitting
		} else if code == 1 {
			stopAt = entity.Cancel
		} else {
			stopAt = entity.Paused
		}
	case <-loopWGChan:
		if task.eofFlag {
			stopAt = entity.Finish
		} else {
			stopAt = entity.Fail
		}
	}
	// signal to sendloop, wait for exit
	close(closeWait)
	loopWG.Wait()
	// close respChan, wait postSend exit
	close(respChan)
	postWG.Wait()

	task.updateProgress()
	task.guard.GuardClose(stopAt)
	task.db.Model(&task.fm).Update("status", stopAt)
	meta, _ := task.fm.Info()
	logger.Info("File", meta["path"], "close - at", stopAt)

	task.reportQuit <- task.Id
}

func (task *Task) sendloop(respChan chan map[string]interface{}, closeWait chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
outer:
	for {
		select {
		case <-closeWait:
			break outer
		default:
			line, fielErr := task.guard.ReadLine()
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				// send url
				if statusCode, err := task.sendUrl(line); err != nil {
					respChan <- map[string]interface{}{"url": line, "statusCode": statusCode, "err": err.Error()}
				} else if statusCode != http.StatusOK {
					respChan <- map[string]interface{}{"url": line, "statusCode": statusCode}
				} else {
					respChan <- nil
				}
			}
			if fielErr != nil {
				// caused by io.EOF or other error
				if fielErr == io.EOF {
					task.eofFlag = true
				}
				break outer
			}
		}
	}
}

func (task *Task) postSend(wg *sync.WaitGroup, respChan chan map[string]interface{}) {
	defer wg.Done()
	var logSize int64 = 0
	var maxLogSize int64 = 20 * 1024 * 1024 // maxLogSize 20M
	logFile, err := os.OpenFile(task.logPath+task.Id+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	var writer *bufio.Writer = nil
	if err == nil {
		if fileInfo, err := logFile.Stat(); err == nil {
			logSize = fileInfo.Size()
		}
		defer logFile.Close()
		writer = bufio.NewWriter(logFile)
	}
outer:
	for i := 0; ; i++ {
		select {
		case res, ok := <-respChan:
			if !ok {
				break outer
			}
			if res == nil {
				task.success++
			} else {
				if writer != nil && logSize < maxLogSize {
					var msg string
					if errMsg, ok := res["err"]; ok {
						msg = fmt.Sprintf("%s\t%d\t%s\n", res["url"], res["statusCode"], errMsg)
					} else {
						msg = fmt.Sprintf("%s\t%d\n", res["url"], res["statusCode"])
					}
					if n, err := writer.WriteString(msg); err == nil {
						logSize += int64(n)
					}
				}
				task.fail++
			}
			if i%500 == 0 {
				task.updateProgress()
			}
		}
	}
	if writer != nil {
		writer.Flush()
	}
}

func (task *Task) sendUrl(url string) (statusCode int, e error) {
	jsonValue, _ := json.Marshal(map[string]interface{}{"url": url})
	for i := 0; i < 3; i++ {
		resp, err := http.Post(task.reciver, "application/json", bytes.NewBuffer(jsonValue))
		// An error is returned if the Client’s CheckRedirect function fails or if there was an HTTP protocol error. A non-2xx response doesn’t cause an error.
		if err == nil {
			defer resp.Body.Close()

			statusCode = resp.StatusCode
			if statusCode != http.StatusOK && statusCode != http.StatusNotFound {
				var resJson interface{}
				if err := json.NewDecoder(resp.Body).Decode(&resJson); err == nil {
					if resBytes, err := json.Marshal(resJson); err == nil {
						e = errors.New(string(resBytes))
					}
				}
			}
			break
		}
		e = err
		time.Sleep(time.Second)
	}
	logger.Debug("send:", url, statusCode)
	return
}

func (task *Task) updateProgress() {
	task.lock.Lock()
	defer task.lock.Unlock()
	fileModel := &entity.FileModel{
		Offset:   task.guard.GuardOffset(),
		Success:  task.success,
		Fail:     task.fail,
		UpdateAt: time.Now(),
	}
	task.db.Model(task.fm).Updates(fileModel)
}

func (task *Task) Info() (info sth.Result, err error) {
	if info, err = task.fm.Info(); err == nil {
		info["offset"] = task.guard.GuardOffset()
		info["success"] = task.success
		info["fail"] = task.fail
	}
	return
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
