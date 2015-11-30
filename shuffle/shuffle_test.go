package shuffle_test

import (
	"fmt"
	"sort"

	"github.com/carlmjohnson/go-utils/shuffle"
)

func Example() {
	//Using a sort.IntSlice since it has the Len and Swap methods already
	s := sort.IntSlice{1, 2, 3, 4}
	shuffle.Shuffle(s)
	fmt.Println(s)
}
