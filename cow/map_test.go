package cow_test

import (
	"testing"

	"github.com/carlmjohnson/go-utils/cow"
)

var m = map[string]string{}

func init() {
	s := "a"
	for i := 0; i < 1024; i++ {
		m[s] = s
		s += "a"
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cow.New(nil)
	}
}

func BenchmarkPlainRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = m["a"]
	}
}

func BenchmarkRead(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		c.Get("a")
	}
}

func BenchmarkInsert(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		c.Insert("a", "a")
	}
}

func BenchmarkUpdate(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		c.Update(m)
	}
}

func BenchmarkRemove(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		c.Remove("x")
	}
}

func BenchmarkReset(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		c.Reset(nil)
		c.Insert("a", "a")
	}
}

func BenchmarkGoRead(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		go c.Get("a")
		c.Get("a")
	}
}

func BenchmarkGoInsert(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		go c.Insert("a", "a")
		c.Insert("a", "a")
	}
}

func BenchmarkGoUpdate(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		go c.Update(m)
		c.Update(m)
	}
}

func BenchmarkGoRemove(b *testing.B) {
	c := cow.New(m)
	for i := 0; i < b.N; i++ {
		go c.Remove("x")
		c.Remove("x")
	}
}
