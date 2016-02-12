package singletasking_test

import (
	"testing"

	"github.com/carlmjohnson/go-utils/singletasking"
)

func blocker() {
	ch := make(chan bool)
	<-ch
}

func void() {}

func BenchmarkBusy(b *testing.B) {
	r := singletasking.New()
	r.Run(blocker)
	for i := 0; i < b.N; i++ {
		r.Run(blocker)
	}
}

func BenchmarkLBusy(b *testing.B) {
	r := singletasking.LRunner{}
	r.Run(blocker)
	for i := 0; i < b.N; i++ {
		r.Run(blocker)
	}
}

func BenchmarkUnblocked(b *testing.B) {
	r := singletasking.New()
	for i := 0; i < b.N; i++ {
		r.Run(void)
	}
}

func BenchmarkLUnblocked(b *testing.B) {
	r := singletasking.LRunner{}
	for i := 0; i < b.N; i++ {
		r.Run(void)
	}
}
