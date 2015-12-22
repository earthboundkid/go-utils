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
	sem  chan struct{}
	done chan struct{}
	// rw ensures that no calls to Release after a call to Stop will
	// result in a token being returned (and thereby acquired by a
	// previously blocked call to Aquire).
	rw sync.RWMutex
	// once ensures that done is not closed twice.
	once sync.Once
}

// New creates a new Semaphore with n max number of tokens allowed.
func New(n int) *Semaphore {
	return &Semaphore{
		sem:  make(chan struct{}, n),
		done: make(chan struct{}),
	}
}

// Acquire returns true after it acquires a token from the underlying
// Semaphore or false if the Semaphore has been closed with Stop().
func (s *Semaphore) Acquire() bool {
	s.rw.RLock()
	sem := s.sem
	s.rw.RUnlock()

	select {
	case sem <- struct{}{}:
		return true
	case <-s.done:
		return false
	}
}

// Release returns a Semaphore token. It is safe to call after the
// Semaphore has been closed with Stop().
func (s *Semaphore) Release() {
	s.rw.RLock()
	sem := s.sem
	s.rw.RUnlock()

	select {
	case <-sem:
	case <-s.done:
	}
}

// Stop closes its underlying Semaphore. It is safe to call multiple
// times.
func (s *Semaphore) Stop() {
	s.once.Do(func() {
		s.rw.Lock()
		defer s.rw.Unlock()

		s.sem = nil
		close(s.done)
	})
}

// Poll reports whether the underlying Semaphore is open. For practical
// purposes, this is only useful for routines that are already holding
// a token from Acquire() if they want to decide to continue working on
// an expensive operation.
func (s *Semaphore) Poll() bool {
	select {
	case <-s.done:
		return false
	default:
		return true
	}
}

func (s *Semaphore) String() string {
	s.rw.RLock()
	capacity, length := cap(s.sem), len(s.sem)
	s.rw.RUnlock()

	return fmt.Sprintf("Semaphore{ n: %d, used: %d }", capacity, length)
}
