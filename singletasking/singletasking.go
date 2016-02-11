// Package singletasking provides a Runner type that runs no more than one task at a time
package singletasking

// Runner will run exactly one task at a time.
type Runner struct {
	run  chan busywork
	stop chan bool
}

type busywork struct {
	task func()
	busy chan bool
}

// New creates a new Runner. Call method Stop when done to prevent resource leak.
func New() Runner {
	r := Runner{
		run:  make(chan busywork),
		stop: make(chan bool),
	}
	go r.init()
	return r
}

func (r Runner) init() {
	busy := false
	done := make(chan bool)
	// Prevent stranding of goroutine if task finishes after Stop
	defer close(done)

	for {
		select {
		case bw := <-r.run:
			bw.busy <- busy
			if !busy {
				go func() {
					bw.task()
					<-done
				}()
				busy = true
			}
		case done <- true:
			busy = false
		case <-r.stop:
			return
		}
	}
}

// Run accepts a task. If the runner is not busy, it returns true and the task is run in a goroutine. If the runner is busy, it returns false and the task is ignored.
func (r Runner) Run(task func()) bool {
	busy := make(chan bool)
	r.run <- busywork{
		task: task,
		busy: busy,
	}
	return !<-busy
}

// Stop causes the Runner to stop listening for tasks.
func (r Runner) Stop() {
	r.stop <- true
}
