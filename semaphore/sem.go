package semaphore

import "sync"

type Semaphore struct {
	sem  chan struct{}
	done chan struct{}
	once sync.Once
}

func New(n int) *Semaphore {
	return &Semaphore{
		sem:  make(chan struct{}, n),
		done: make(chan struct{}),
	}
}

func (s *Semaphore) Acquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	case <-s.done:
		return false
	}
}

func (s *Semaphore) Release() {
	go func() {
		select {
		case <-s.sem:
		case <-s.done:
		}
	}()
}

func (s *Semaphore) Stop() {
	s.once.Do(func() {
		s.sem = nil
		close(s.done)
	})
}
