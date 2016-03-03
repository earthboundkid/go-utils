package semaphore_test

import (
	"sync"
	"testing"
	"time"

	"github.com/carlmjohnson/go-utils/semaphore"
)

func BenchmarkSemAcq(b *testing.B) {
	s := semaphore.New(1)
	for i := 0; i < b.N; i++ {
		s.Acquire()
		s.Release()
	}
}

func BenchmarkLSemAcq(b *testing.B) {
	s := semaphore.NewLS(1)
	for i := 0; i < b.N; i++ {
		s.Acquire()
		s.Release()
	}
}

func BenchmarkSemAcqRelease(b *testing.B) {
	s := semaphore.New(1)
	s.Acquire()

	for i := 0; i < b.N; i++ {
		go s.Release()
		s.Acquire()
	}
}

func BenchmarkLSemAcqRelease(b *testing.B) {
	s := semaphore.NewLS(1)
	s.Acquire()

	for i := 0; i < b.N; i++ {
		go s.Release()
		s.Acquire()
	}
}

func BenchmarkSemAcqStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.New(1)
		s.Acquire()
		go s.Stop()
		s.Acquire()
	}
}

func BenchmarkLSemAcqStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.NewLS(1)
		s.Acquire()
		go s.Stop()
		s.Acquire()
	}
}

func BenchmarkSemLoop1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.New(1)
		go s.Stop()
		for s.Acquire() {
		}
	}
}

func BenchmarkLSemLoop1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.NewLS(1)
		go s.Stop()
		for s.Acquire() {
		}
	}
}

func BenchmarkSemLoop5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.New(5)
		go s.Stop()
		for s.Acquire() {
		}
	}
}

func BenchmarkLSemLoop5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.NewLS(5)
		go s.Stop()
		for s.Acquire() {
		}
	}
}

func BenchmarkSemLoop50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.New(50)
		go s.Stop()
		for s.Acquire() {
		}
	}
}

func BenchmarkLSemLoop50(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.NewLS(50)
		go s.Stop()
		for s.Acquire() {
		}
	}
}

func BenchmarkSemSleep(b *testing.B) {
	s := semaphore.New(1)
	s.Acquire()
	for i := 0; i < b.N; i++ {
		go func() {
			time.Sleep(5 * time.Nanosecond)
			s.Release()
		}()
		s.Acquire()
	}
}

func BenchmarkLSemSleep(b *testing.B) {
	s := semaphore.NewLS(1)
	s.Acquire()
	for i := 0; i < b.N; i++ {
		go func() {
			time.Sleep(5 * time.Nanosecond)
			s.Release()
		}()
		s.Acquire()
	}
}

func BenchmarkSemWait(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.New(10)
		for j := 0; j < 100; j++ {
			go func() {
				for s.Acquire() {
					time.Sleep(1 * time.Nanosecond)
					s.Release()
				}
			}()
		}
		time.Sleep(100 * time.Nanosecond)
		s.Stop()
		s.Wait()
	}
}

func BenchmarkLSemWait(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := semaphore.NewLS(10)

		var wg sync.WaitGroup
		wg.Add(100)
		for j := 0; j < 100; j++ {
			go func() {
				defer wg.Done()
				for s.Acquire() {
					time.Sleep(1 * time.Nanosecond)
					s.Release()
				}
			}()
		}
		time.Sleep(100 * time.Nanosecond)
		s.Stop()
		wg.Wait()
	}
}
