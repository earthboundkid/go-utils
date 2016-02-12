package singletasking

import "sync"

// LRunner will run exactly one task at a time.
//
// Unlike Runner, it uses locking internally and does not need to be stopped to prevent resource leak.
type LRunner struct {
	m    sync.Mutex
	busy bool
}

// Run accepts a task. If the runner is not busy, it returns true and the task is run in a goroutine. If the runner is busy, it returns false and the task is ignored.
func (r *LRunner) Run(task func()) bool {
	r.m.Lock()
	defer r.m.Unlock()
	run := !r.busy
	if run {
		r.busy = true
		go func() {
			task()
			r.m.Lock()
			defer r.m.Unlock()
			r.busy = false
		}()
	}
	return run
}
