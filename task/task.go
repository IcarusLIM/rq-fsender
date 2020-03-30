package task

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

// ProcessStatus type
type ProcessStatus string

// ProcessStatus enum
const (
	Waitting   ProcessStatus = "waitting"
	Finished   ProcessStatus = "finished"
	Canceled   ProcessStatus = "canceled"
	Processing ProcessStatus = "processing"
)

// Task TODO
type Task struct {
	id         string
	tc         *Container
	quit       *chan int
	SC         *chan *Feed
	Receiver   string
	Stats      *Stats
	Files      []PlainFile
	Index      int
	CurrReader VisableReader
}

// NewTask TODO
func NewTask(taskID string, tc *Container, sc *chan *Feed, receiver string, files []PlainFile) Task {
	quit := make(chan int, 1)
	return Task{
		id:         taskID,
		tc:         tc,
		quit:       &quit,
		SC:         sc,
		Receiver:   receiver,
		Stats:      NewStats(),
		Files:      files,
		Index:      0,
		CurrReader: nil,
	}
}

// Run task
func (t *Task) Run() {
	for ; t.Index < len(t.Files); t.Index++ {
		isCanceled := t.parse(t.Files[t.Index])
		if isCanceled {
			break
		}
	}
	t.tc.Clean(t.id)
}

func (t *Task) parse(file PlainFile) bool {
	var err error
	t.CurrReader, err = file.Open()
	if err != nil {
		file.SetFlag(&Result{"process": "error", "err": err.Error()})
		return false
	}
	reader := bufio.NewReader(t.CurrReader)
	for {
		time.Sleep(time.Second)
		select {
		case <-*t.quit:
			return true
		default:
			line, _, err := reader.ReadLine()
			if err != nil {
				return false
			}
			fmt.Print("Read line:" + string(line) + "\n")
			req := Request{
				URL:    string(line),
				Method: "GET",
			}
			*t.SC <- &Feed{
				Receiver: t.Receiver,
				Req:      req,
				counter:  t.Stats,
			}
		}
	}
}

// Status of task
func (t *Task) Status() Result {
	res := Result{}
	res["file_num"] = len(t.Files)
	res["file_cur"] = t.Index
	currentReader := t.CurrReader
	if currentReader != nil {
		size, _ := currentReader.Size()
		offset, _ := currentReader.Seek(0, io.SeekCurrent)
		res["cur_size"] = size
		res["cur_offset"] = offset
	}
	res["stats"] = t.Stats.Stats()
	return res
}

// LogT TODO
func (t *Task) LogT() Result {
	res := Result{}
	res["file_num"] = len(t.Files)
	res["stats"] = t.Stats.Report()
	return res
}

// Container TODO
type Container struct {
	tasks map[string]Task
	lock  *sync.Mutex
}

// NewContainer TODO
func NewContainer() *Container {
	tasks := make(map[string]Task)
	return &Container{tasks: tasks, lock: &sync.Mutex{}}
}

// Add TODO
func (tc *Container) Add(sc *chan *Feed, taskID string, receiver string, files []PlainFile) error {
	tc.lock.Lock()
	defer tc.lock.Unlock()
	if _, ok := tc.tasks[taskID]; ok {
		return errors.New("Dup taskID")
	}
	task := NewTask(taskID, tc, sc, receiver, files)
	tc.tasks[taskID] = task
	go task.Run()
	return nil
}

// Get Task
func (tc *Container) Get(taskID string) (task Task, ok bool) {
	task, ok = tc.tasks[taskID]
	return
}

// List Task
func (tc *Container) List() ([]string, error) {
	fmt.Println(len(tc.tasks))
	keys := make([]string, 0, len(tc.tasks))
	for k := range tc.tasks {
		keys = append(keys, k)
	}
	return keys, nil
}

// Cancel Task
func (tc *Container) Cancel(taskID string) (Result, error) {
	task, ok := tc.Get(taskID)
	if !ok {
		return nil, errors.New("Task not exit")
	}
	delete(tc.tasks, taskID)
	*(task.quit) <- 0
	close(*(task.quit))
	return task.LogT(), nil
}

// Clean task
func (tc *Container) Clean(taskID string) {
	task, ok := tc.Get(taskID)
	if !ok {
		return
	}
	delete(tc.tasks, taskID)
	fmt.Print(task.LogT())
}
