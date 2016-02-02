package semaphore

import (
	"fmt"
	"sync"
)

// LockingSemaphore is functionally identical to Semaphore, but uses
// mutexes internally.
type LockingSemaphore struct {
	sem  chan struct{}
	done chan struct{}
	// rw ensures that no calls to Release after a call to Stop will
	// result in a token being returned (and thereby acquired by a
	// previously blocked call to Aquire).
	rw sync.RWMutex
	// once ensures that done is not closed twice.
	once sync.Once
}

// NewLS creates a new LockingSemaphore with n max number of tokens allowed.
func NewLS(n int) *LockingSemaphore {
	return &LockingSemaphore{
		sem:  make(chan struct{}, n),
		done: make(chan struct{}),
	}
}

// Acquire returns true after it acquires a token from the underlying
// LockingSemaphore or false if the LockingSemaphore has been closed
// with Stop().
func (s *LockingSemaphore) Acquire() bool {
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

// Release returns a LockingSemaphore token. It is safe to call after
// the LockingSemaphore has been closed with Stop().
func (s *LockingSemaphore) Release() {
	s.rw.RLock()
	sem := s.sem
	s.rw.RUnlock()

	select {
	case <-sem:
	case <-s.done:
	}
}

// Stop closes its underlying LockingSemaphore. It is safe to call multiple
// times.
func (s *LockingSemaphore) Stop() {
	s.once.Do(func() {
		s.rw.Lock()
		defer s.rw.Unlock()

		s.sem = nil
		close(s.done)
	})
}

// Poll reports whether the underlying LockingSemaphore is open. For practical
// purposes, this is only useful for routines that are already holding
// a token from Acquire() if they want to decide to continue working on
// an expensive operation.
func (s *LockingSemaphore) Poll() bool {
	select {
	case <-s.done:
		return false
	default:
		return true
	}
}

func (s *LockingSemaphore) String() string {
	s.rw.RLock()
	capacity, length := cap(s.sem), len(s.sem)
	s.rw.RUnlock()

	return fmt.Sprintf("LockingSemaphore{ n: %d, used: %d }", capacity, length)
}
