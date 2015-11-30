shuffle
=======

Implementation of the [Fisher–Yates shuffle](http://en.wikipedia.org/wiki/Fisher–Yates_shuffle) (or Knuth shuffle) in Go.

Usage
-----

	//Using a sort.IntSlice since it has the Len and Swap methods already
	s := sort.IntSlice{1, 2, 3, 4}
	shuffle.Shuffle(s)
	fmt.Println(s)


