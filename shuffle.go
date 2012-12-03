// By Carl M. Johnson, MIT License

// Package shuffle implements the Fisher–Yates shuffle (or Knuth shuffle).
//
// See http://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
package shuffle

import "math/rand"

//Type Interface is similar to a sort.Interface, but there's no Less method
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

// Shuffle shuffles the data using the following algorithm:
//   To shuffle an array a of n elements (indices 0..n-1):
//     for i from n − 1 downto 1 do
//       j ← random integer with 0 ≤ j ≤ i
//       exchange a[j] and a[i]
// (thanks Wikipedia)
func Shuffle(s Interface) {
	for i := s.Len() - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s.Swap(i, j)
	}
}
