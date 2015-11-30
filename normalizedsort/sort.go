// Package normalizedsort contains utilities for normalized sorting of []string.
package normalizedsort

import (
	"sort"
	"strings"
)

// Sort is a convenience method that calls New then Go's sort.Sort on the resulting sort.Interface.
func Sort(ss []string, normalize func(string) string) {
	sortable := New(ss, normalize)
	sort.Sort(sortable)
}

// New returns a sort.Interface that sorts according to its normalization function. If normalize is nil, the resulting sort.Interface uses strings.ToLower by default. If two strings normalize to the same value, the interface sorts them according to their unnormalized form, i.e. upper-case comes before lower-case.
func New(ss []string, normalize func(string) string) sort.Interface {
	if normalize == nil {
		normalize = strings.ToLower
	}
	sortable := normalizedStringSlice{
		original: ss,
	}
	sortable.init(normalize)
	return &sortable
}

type normalizedStringSlice struct {
	original   []string
	normalized []string
}

func (ns *normalizedStringSlice) init(normalize func(string) string) {
	ns.normalized = make([]string, 0, len(ns.original))
	for i := range ns.original {
		ns.normalized = append(ns.normalized, normalize(ns.original[i]))
	}
}

// Len is the number of elements in the collection.
func (ns *normalizedStringSlice) Len() int {
	return len(ns.original)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (ns *normalizedStringSlice) Less(i, j int) bool {
	// If there's a tie, use the original strings to sort
	if ns.normalized[i] == ns.normalized[j] {
		return ns.original[i] < ns.original[j]
	}
	return ns.normalized[i] < ns.normalized[j]
}

// Swap swaps the elements with indexes i and j.
func (ns *normalizedStringSlice) Swap(i, j int) {
	ns.original[i], ns.original[j] = ns.original[j], ns.original[i]
	ns.normalized[i], ns.normalized[j] = ns.normalized[j], ns.normalized[i]
}

// CaseInsensitiveTrimSpace calls strings.TrimSpace and strings.ToLower to normalize its input.
func CaseInsensitiveTrimSpace(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
