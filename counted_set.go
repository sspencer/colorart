package colorart

import (
	"fmt"
	"sort"
)

const maxComponent = 255

// rgb is RGB components, 0-255, used as the map key
type rgb [3]byte

// CountedSet counts the number of times each object (string) is added to the set.
// The set is not thread safe.
type CountedSet map[rgb]int

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

//---------------------------

// NewCountedSet creates a new CountedSet of the specified size.
func NewCountedSet(size int) CountedSet {
	s := make(map[rgb]int, size)
	return s
}

// AddPixel converts pixel to [3]byte rgb and counts unique colors
func (s CountedSet) AddPixel(p pixel) {

	b := uint8(colorShifter)
	ri := uint8(maxComponent*p.R) >> b << b
	gi := uint8(maxComponent*p.G) >> b << b
	bi := uint8(maxComponent*p.B) >> b << b

	color := rgb{ri, gi, bi}

	s[color]++
}

// Merge other counted set into this one.
func (s CountedSet) Merge(o CountedSet) {
	for color, cnt := range o {
		s[color] += cnt
	}
}

// Add color with count.
func (s CountedSet) AddCount(color rgb, count int) {
	s[color] += count
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
