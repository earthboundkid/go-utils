// Package singletasking provides a Runner type that runs no more than one task at a time
package singletasking

// Runner will run exactly one task at a time.
type Runner struct {
	runtask chan func()
	runbusy chan bool
	stop    chan bool
}

// New creates a new Runner. Call method Stop when done to prevent resource leak.
func New() Runner {
	r := Runner{
		runtask: make(chan func()),
		runbusy: make(chan bool),
		stop:    make(chan bool),
	}
	go r.init()
	return r
}

func (r Runner) init() {
	var busyCh chan bool
	taskCh := r.runtask
	done := make(chan bool)
	// Prevent stranding of goroutine if task finishes after Stop
	defer close(done)
	defer close(r.runbusy)
	defer close(r.stop)

	for {
		select {
		case busyCh <- true:

		case task := <-taskCh:
			go func() {
				task()
				<-done
			}()
			busyCh = r.runbusy
			taskCh = nil

		case done <- true:
			busyCh = nil
			taskCh = r.runtask

		case r.stop <- true:
			return
		}
	}
}

// Run accepts a task. If the runner is not busy, it returns true and the task is run in a goroutine. If the runner is busy, it returns false and the task is ignored.
func (r Runner) Run(task func()) bool {
	select {
	case r.runtask <- task:
		return true
	case <-r.runbusy:
		return false
	}
}

// Stop causes the Runner to stop listening for tasks.
func (r Runner) Stop() {
	<-r.stop
}
