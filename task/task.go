package task

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Ghamster0/os-rq-fsender/sender"
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
	SC         *chan *sender.Feed
	Receiver   string
	Files      []PlainFile
	Index      int
	CurrReader VisableReader
}

// Run task
func (t *Task) Run() {
	t.Index = 0
	processing := t.Files[t.Index]
	t.CurrReader, _ = processing.Open()
	reader := bufio.NewReader(t.CurrReader)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		fmt.Print("Read line:" + string(line) + "\n")
		req := sender.Request{
			URL:    string(line),
			Method: "GET",
		}
		*t.SC <- &sender.Feed{
			Receiver: t.Receiver,
			Req:      req,
		}
	}
}

// Status of task
func (t *Task) Status() (res Result) {
	res["pedding"] = 1
	currReader := t.CurrReader
	size, _ := currReader.Size()
	pos, _ := currReader.Seek(0, io.SeekCurrent)
	res["process"] = Result{"size": size, "pos": pos}
	return
}

// Container TODO
type Container struct {
	tasks map[string]Task
}

// NewContainer TODO
func NewContainer() *Container {
	tasks := make(map[string]Task)
	return &Container{tasks: tasks}
}

// AddTask TODO
func (tc *Container) AddTask(sc *chan *sender.Feed, taskID string, receiver string, files []PlainFile) {
	task := Task{
		SC:         sc,
		Receiver:   receiver,
		Files:      files,
		Index:      0,
		CurrReader: nil,
	}
	go task.Run()
	tc.tasks[taskID] = task
	return
}
