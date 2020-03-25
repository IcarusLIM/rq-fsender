package sender

import "fmt"

// Runner TODO
type Runner struct {
	SC chan Feed
}

// Start workers
func (runner *Runner) Start(threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			fc := runner.SC
			for feed := range fc {
				fmt.Print(feed.Receiver)
				feed.send()
			}
		}()
	}
}
