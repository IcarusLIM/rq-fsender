package task

import (
	"sync/atomic"
	"time"

	"github.com/prep/average"
)

// Stats TODO
type Stats struct {
	total     int64
	success   int64
	fail      int64
	successWS *average.SlidingWindow
	failWS    *average.SlidingWindow
}

// Stats TODO
func (t *Stats) Stats() Result {
	success5s, _ := t.successWS.Total(5 * time.Second)
	fail5s, _ := t.failWS.Total(5 * time.Second)
	return Result{"total": t.total, "success": t.success, "fail": t.fail, "success_5s": success5s, "fail_5s": fail5s}
}

// Report TODO
func (t *Stats) Report() Result {
	return Result{"total": t.total, "success": t.success, "fail": t.fail}
}

// OnSuccess TODO
func (t *Stats) OnSuccess() {
	atomic.AddInt64(&(t.total), 1)
	atomic.AddInt64(&(t.success), 1)
	t.successWS.Add(1)
}

// OnFail TODO
func (t *Stats) OnFail() {
	atomic.AddInt64(&(t.total), 1)
	atomic.AddInt64(&(t.fail), 1)
	t.failWS.Add(1)
}

// NewStats TODO
func NewStats() *Stats {
	return &Stats{
		0,
		0,
		0,
		average.MustNew(time.Minute, time.Second),
		average.MustNew(time.Minute, time.Second),
	}
}
