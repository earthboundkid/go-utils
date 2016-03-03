// Package semaphore is a package containing a simple reference
// implementation of a semaphore type. If you need to use a semaphore in
// your own package, probably the right way to use this package is to
// copy it and make your own specialized implementation.
package semaphore

import "fmt"

// A Semaphore helps add restrictions on the number of active goroutines
// by ensuring that no more than n tokens may be acquired at one time.
// It also allows for the broadcasting of Stop messages.
type Semaphore struct {
	// use chan bool to simplify testing if channel is closed
	acquire  chan bool
	release  chan struct{}
	stop     chan struct{}
	wait     chan struct{}
	poll     chan struct{}
	stringer chan stringerData
}

type stringerData struct {
	max, count int
}

// New creates a new Semaphore with n max number of tokens allowed.
func New(n int) *Semaphore {
	s := Semaphore{
		acquire:  make(chan bool),
		release:  make(chan struct{}),
		stop:     make(chan struct{}),
		wait:     make(chan struct{}),
		poll:     make(chan struct{}),
		stringer: make(chan stringerData),
	}
	go s.start(n)
	return &s
}

func (s *Semaphore) start(max int) {
	count := 0

	for {
		var acquire = s.acquire

		// nil always blocks sends
		if count >= max {
			acquire = nil
		}

		select {
		case acquire <- true:
			count++
		case s.release <- struct{}{}:
			count--
		case s.stringer <- stringerData{max, count}:
		case s.stop <- struct{}{}:
			close(s.acquire)
			close(s.stringer)
			close(s.poll)
			close(s.stop)

			for count > 0 {
				s.release <- struct{}{}
				count--
			}

			close(s.wait)
			return
		}
	}
}

// Acquire returns true after it acquires a token from the underlying
// Semaphore or false if the Semaphore has been closed with Stop().
func (s *Semaphore) Acquire() bool {
	return <-s.acquire
}

// Release returns a Semaphore token. It is safe to call after the
// Semaphore has been closed with Stop().
func (s *Semaphore) Release() {
	<-s.release
}

// Stop closes its underlying Semaphore. It is safe to call multiple
// times.
func (s *Semaphore) Stop() {
	<-s.stop
}

// Wait blocks until all acquired tokens have been released.
func (s *Semaphore) Wait() {
	<-s.wait
}

// Poll reports whether the underlying Semaphore is open. For practical
// purposes, this is only useful for routines that are already holding
// a token from Acquire() if they want to decide to continue working on
// an expensive operation.
func (s *Semaphore) Poll() bool {
	select {
	case <-s.poll:
		return false
	default:
		return true
	}
}

func (s *Semaphore) String() string {
	if v, ok := <-s.stringer; ok {
		return fmt.Sprintf("Semaphore{ n: %d, used: %d }", v.max, v.count)
	}
	return "Semaphore{closed}"
}
