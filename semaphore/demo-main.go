// Allow package to be included in the same directory:
// +build ignore

// Demonstration of Semaphore package. No more than N routines should be
// running at any one time.
package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/carlmjohnson/go-utils/semaphore"
)

// N is the number of active semaphores allowed.
const N = 3

var (
	s = semaphore.New(N)
	w = new(int32)
)

func worker(n int) {
	for s.Acquire() {
		fmt.Printf("%s\t%d\t%*s\n",
			time.Now().Format("15:04:05"), atomic.AddInt32(w, 1), n, "^")
		time.Sleep(1 * time.Second)
		fmt.Printf("%s\t%d\t%*s\n",
			time.Now().Format("15:04:05"), atomic.AddInt32(w, -1), n, "_")
		s.Release()
	}
	fmt.Printf("%s\t%d\t%*s\n", time.Now().Format("15:04:05"), atomic.LoadInt32(w), n, "*")
}

func main() {
	fmt.Println("Time\t# of workers\n")
	for i := 0; i < 10; i++ {
		go worker(5 * i)
	}
	go func() {
		time.Sleep(10 * time.Second)
		s.Stop()
	}()
	s.Wait()
}
