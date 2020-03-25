package ctx

import (
	"github.com/Ghamster0/os-rq-fsender/sender"
	"github.com/Ghamster0/os-rq-fsender/task"
)

// ApplicationContext TODO
type ApplicationContext struct {
	SC    chan sender.Feed
	Tasks *task.TaskContainer
}

// StartApplication TODO
func StartApplication() (ctx *ApplicationContext, err error) {
	threads := 10
	sc := make(chan sender.Feed, 2*threads)
	runner := &sender.Runner{SC: sc}
	runner.Start(threads)
	ctx = &ApplicationContext{
		SC: sc,
	}
	return
}
