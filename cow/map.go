// Package cow provides a reference implementation of copy-on-write maps
// for read-heavy, write-light data.
package cow

import (
	"sync"
	"sync/atomic"
)

// Map is a synchronous copy on write map. Reads are cheap. Writes are expensive.
type Map struct {
	v atomic.Value
	m sync.Mutex // used only by writers
}

func copyM(dst map[string]string, src map[string]string) {
	for k, v := range src {
		dst[k] = v // copy all data from the current object to the new one
	}
}

func dup(src map[string]string) (dst map[string]string) {
	dst = make(map[string]string, len(src))
	copyM(dst, src)
	return dst
}

// New initializes a new map based on m.
func New(m map[string]string) *Map {
	cowm := Map{}
	cowm.v.Store(dup(m))
	return &cowm
}

// Get retreives the value associated with the key from the Map.
func (cowm *Map) Get(key string) (val string) {
	// No lock needed!
	m := cowm.v.Load().(map[string]string)
	return m[key]
}

// Insert inserts a key-value pair.
func (cowm *Map) Insert(key, val string) {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	src := cowm.v.Load().(map[string]string)
	dst := dup(src)
	dst[key] = val
	cowm.v.Store(dst)
}

// Update efficiently inserts all the values in m into the Map.
func (cowm *Map) Update(m map[string]string) {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	src := cowm.v.Load().(map[string]string)
	dst := dup(src)
	copyM(dst, m)
	cowm.v.Store(dst)
}

// Remove removes key from the Map.
func (cowm *Map) Remove(key string) {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	src := cowm.v.Load().(map[string]string)
	dst := dup(src)
	delete(dst, key)
	cowm.v.Store(dst)
}

// Reset initializes the Map to the values in m. Use of nil to empty the Map is okay.
func (cowm *Map) Reset(m map[string]string) {
	cowm.m.Lock()
	defer cowm.m.Unlock()

	cowm.v.Store(dup(m))
}
