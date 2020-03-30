package sender

import "github.com/Ghamster0/os-rq-fsender/task"

// Runner TODO
type Runner struct {
	// Sender channel
	SC *chan *task.Feed
}

// Start workers
func (runner *Runner) Start(threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			fc := runner.SC
			for feed := range *fc {
				feed.Send()
			}
		}()
	}
}
