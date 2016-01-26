// Package semaphore is a package containing a simple reference
// implementation of a semaphore type. If you need to use a semaphore in
// your own package, probably the right way to use this package is to
// copy it and make your own specialized implementation.
package semaphore

import (
	"fmt"
	"sync"
)

// A Semaphore helps add restrictions on the number of active goroutines
// by ensuring that no more than n tokens may be acquired at one time.
// It also allows for the broadcasting of Stop messages.
type Semaphore struct {
	max         int
	count       int
	acquire     chan struct{}
	doneAcquire chan struct{}
	release     chan struct{}
	doneRelease chan struct{}
	stop        chan struct{}
	// once ensures that stop/done are not closed twice.
	once sync.Once
}

// New creates a new Semaphore with n max number of tokens allowed.
func New(n int) *Semaphore {
	s := Semaphore{
		max:         n,
		count:       0,
		acquire:     make(chan struct{}),
		doneAcquire: make(chan struct{}),
		release:     make(chan struct{}),
		doneRelease: make(chan struct{}),
		stop:        make(chan struct{}),
	}
	go s.start()
	return &s
}

func (s *Semaphore) start() {
	for {
		if s.count < s.max {
			select {
			case <-s.release:
				s.count--
			case <-s.acquire:
				s.count++
			case <-s.stop:
				close(s.doneAcquire)
				return
			}
		} else {
			select {
			case <-s.release:
				s.count--
			case <-s.stop:
				close(s.doneAcquire)
				return
			}
		}
	}
}

// Acquire returns true after it acquires a token from the underlying
// Semaphore or false if the Semaphore has been closed with Stop().
func (s *Semaphore) Acquire() bool {
	select {
	case s.acquire <- struct{}{}:
		return true
	case <-s.doneAcquire:
		return false
	}
}

// Release returns a Semaphore token. It is safe to call after the
// Semaphore has been closed with Stop().
func (s *Semaphore) Release() {
	select {
	case s.release <- struct{}{}:
	case <-s.doneRelease:
	}
}

// Stop closes its underlying Semaphore. It is safe to call multiple
// times. If wait is true, it will block until all semaphores are
// released.
func (s *Semaphore) Stop(wait bool) {
	s.once.Do(func() {
		if wait {
			// Block until .start has returned...
			s.stop <- struct{}{}
			// Then drain the channel
			for s.count > 0 {
				<-s.release
				s.count--
			}
			// Shouldn't matter unless there's a programmatic error, but
			// let's close this anyway...
			close(s.doneRelease)
		} else {
			close(s.stop)
			close(s.doneRelease)
		}
	})
}

// Poll reports whether the underlying Semaphore is open. For practical
// purposes, this is only useful for routines that are already holding
// a token from Acquire() if they want to decide to continue working on
// an expensive operation.
func (s *Semaphore) Poll() bool {
	select {
	case <-s.doneAcquire:
		return false
	default:
		return true
	}
}

// Warning: This is racy
func (s *Semaphore) String() string {
	return fmt.Sprintf("Semaphore{ n: %d, used: %d }", s.max, s.count)
}
