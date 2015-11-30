package normalizedsort_test

import (
	"fmt"

	"github.com/carlmjohnson/go-utils/normalizedsort"
)

func ExampleSort() {
	slice := []string{"Aardvark", "hello", "aardvark", "  Hello", "World!"}
	normalizedsort.Sort(slice, nil)
	fmt.Printf("%q\n", slice)
	// Output: ["  Hello" "Aardvark" "aardvark" "hello" "World!"]
}

func ExampleCaseInsensitiveTrimSpace() {
	slice := []string{"Aardvark", "hello", "aardvark", "  Hello", "World!"}
	normalizedsort.Sort(slice, normalizedsort.CaseInsensitiveTrimSpace)
	fmt.Printf("%q\n", slice)
	// Output: ["Aardvark" "aardvark" "  Hello" "hello" "World!"]
}
