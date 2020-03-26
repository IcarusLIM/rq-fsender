package ctx

import (
	"github.com/Ghamster0/os-rq-fsender/sender"
	"github.com/Ghamster0/os-rq-fsender/task"
)

// ApplicationContext TODO
type ApplicationContext struct {
	SC    chan *sender.Feed
	Tasks *task.Container
}

// StartApplication TODO
func StartApplication() (ctx *ApplicationContext, err error) {
	threads := 2
	sc := make(chan *sender.Feed, 2*threads)
	// start sender runner
	runner := &sender.Runner{SC: &sc}
	runner.Start(threads)
	// init taskContainer
	tasks := task.NewContainer()
	ctx = &ApplicationContext{
		SC:    sc,
		Tasks: tasks,
	}
	return
}
