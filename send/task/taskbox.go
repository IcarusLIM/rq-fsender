package task

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var logger, _ = logging.GetLogger("TASK")

type TaskBox struct {
	db         *gorm.DB
	conf       *viper.Viper
	stop       chan int // must be blocking channel
	removing   chan string
	tasks      map[string]*Task
	loc        *sync.RWMutex
	concurrent int
}

func NewTaskBox(db *gorm.DB, conf *viper.Viper) *TaskBox {
	concurr := conf.GetInt("task.concurrent")
	return &TaskBox{
		db:         db,
		conf:       conf,
		stop:       make(chan int),
		removing:   make(chan string, concurr),
		tasks:      make(map[string]*Task),
		loc:        &sync.RWMutex{},
		concurrent: concurr,
	}
}

func TaskBoxServe(lifecycle fx.Lifecycle, box *TaskBox) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				box.OnStart()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				box.OnStop()
				return nil
			},
		},
	)
}

func (box *TaskBox) OnStart() {
	go box.run()
}

func (box *TaskBox) run() {
loop:
	for {
		logger.Debug("Tasks running - ", len(box.tasks))
		select {
		case <-time.After(time.Second):
			box.createTask()
		case id := <-box.removing:
			box.cleanTask(id)
		case <-box.stop:
			break loop
		}
	}
}

func (box *TaskBox) createTask() {
	box.loc.Lock()
	defer box.loc.Unlock()
	if len(box.tasks) >= box.concurrent {
		return
	}
	if task, err := NewTask(box.removing, box.db, box.conf); err == nil {
		if err := task.Run(); err == nil {
			box.tasks[task.Id] = task
		}
	}
}

func (box *TaskBox) cleanTask(id string) {
	box.loc.Lock()
	defer box.loc.Unlock()
	delete(box.tasks, id)
}

func (box *TaskBox) OnStop() {
	// wait for run() exit
	box.stop <- 0
	// then acquire loc
	box.loc.Lock()
	defer box.loc.Unlock()
	for key := range box.tasks {
		task := box.tasks[key]
		task.Stop()
	}
	box.waitTasksExit()
}

func (box *TaskBox) waitTasksExit() {
	if len(box.tasks) > 0 {
		for {
			select {
			case id := <-box.removing:
				delete(box.tasks, id)
			case <-time.After(time.Second):
			}
			if len(box.tasks) == 0 {
				break
			}
		}
	}
}

func (box *TaskBox) InfoTask(key string) (sth.Result, error) {
	task, b := box.tasks[key]
	if b {
		return task.Info()
	} else {
		return nil, errors.New("")
	}
}

func (box *TaskBox) PauseTask(key string) error {
	box.loc.Lock()
	defer box.loc.Unlock()
	task, b := box.tasks[key]
	if b {
		task.Pause()
		return nil
	} else {
		return errors.New("")
	}
}

func (box *TaskBox) CancelTask(key string) (sth.Result, error) {
	box.loc.Lock()
	defer box.loc.Unlock()
	task, b := box.tasks[key]
	if b {
		task.Cancel()
		return task.Info()
	} else {
		return nil, errors.New("")
	}
}
