package task

import (
	"bufio"
	"fmt"

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
	SC       chan *sender.Feed
	Receiver string
	Files    []ReadOnlyFile
	Status   map[ProcessStatus][]ReadOnlyFile
}

// Run task
func (t *Task) Run() {
	metaFile := t.Files[0]
	file, _ := metaFile.Open()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		fmt.Print("extract line: ", line)
		req := sender.Request{
			URL:    string(line),
			Method: "GET",
		}
		t.SC <- &sender.Feed{
			Receiver: t.Receiver,
			Req:      req,
		}
	}
}

// TaskContainer TODO
type TaskContainer struct {
	tasks map[string]Task
}

func (tc *TaskContainer) AddTask() {
	return
}
