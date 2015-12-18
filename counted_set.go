package colorart

import (
	"fmt"
	"sort"
)

type rgb [3]byte

// CountedEntry is for use by sorting class
type CountedEntry struct {
	Color rgb
	Count int
}

// ByCount is the type used to sort
type ByCount []CountedEntry

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].Count > a[j].Count }

func (e CountedEntry) String() string {
	return fmt.Sprintf("%02x%02x%02x: %d", e.Color[0], e.Color[1], e.Color[2], e.Count)
}

// CountedSet counts the number of times each object (string) is added to the set.
// The set is not thread safe.

type CountedSet map[rgb]int

//---------------------------

// NewCountedSet creates a new CountedSet of the specified size.
func NewCountedSet(size int) CountedSet {
	s := make(map[rgb]int, size)
	return s
}

// Add adds an object to the set.
func (s CountedSet) Add(color rgb) {
	s[color]++
}

// AddPixel converts pixel to [3]byte rgb and counts unique colors
func (s CountedSet) AddPixel(p pixel) {
	const max = 255

	/*
		b := uint8(3)
		ri := uint8(max*p.R) >> b << b
		gi := uint8(max*p.G) >> b << b
		bi := uint8(max*p.B) >> b << b
	*/
	ri := uint8(max * p.R)
	gi := uint8(max * p.G)
	bi := uint8(max * p.B)

	color := rgb{ri, gi, bi}

	s[color]++
}

// Merge other counted set into this one.
func (s CountedSet) Merge(o CountedSet) {
	for color, cnt := range o {
		s.AddCount(color, cnt)
	}
}

// Add color with count.
func (s CountedSet) AddCount(color rgb, count int) {
	s[color] = count
}

// Size returns the number of objects in the set.
func (s CountedSet) Size() int {
	return len(s)
}

// Count returns the number of times the specified object has been added to the set.
func (s CountedSet) Count(color rgb) int {
	return s[color]
}

// SortedSet returns the entries (Color, Count) ordered from greatest count to least
func (s CountedSet) SortedSet() []CountedEntry {

	list := make([]CountedEntry, len(s))
	it := 0
	for color, cnt := range s {
		list[it] = CountedEntry{color, cnt}
		it++
	}

	sort.Sort(ByCount(list))
	return list
}
