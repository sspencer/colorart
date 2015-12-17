package colorart

import (
	"fmt"
	"sort"
)

type rgb [3]byte

func RGBAToRGB(r, g, b, a uint32) rgb {
	const max = 255
	fa := float64(a)
	ri := uint8(max * float64(r) / fa)
	gi := uint8(max * float64(g) / fa)
	bi := uint8(max * float64(b) / fa)

	return rgb{ri, gi, bi}
}

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
type CountedSet struct {
	m map[rgb]int
}

//---------------------------

// NewCountedSet creates a new CountedSet of the specified size.
func NewCountedSet(size int) *CountedSet {
	s := &CountedSet{}
	s.m = make(map[rgb]int, size)
	return s
}

// Add adds an object to the set.
func (s *CountedSet) Add(color rgb) {
	s.m[color]++
}

// AddPixel converts pixel to [3]byte rgb and counts unique colors
func (s *CountedSet) AddPixel(p pixel) {
	const max = 255
	ri := uint8(max * p.R)
	gi := uint8(max * p.G)
	bi := uint8(max * p.B)

	color := rgb{ri, gi, bi}
	s.m[color]++
}

func (s *CountedSet) AddCount(color rgb, count int) {
	s.m[color] = count
}

// Size returns the number of objects in the set.
func (s *CountedSet) Size() int {
	return len(s.m)
}

// Count returns the number of times the specified object has been added to the set.
func (s *CountedSet) Count(color rgb) int {
	return s.m[color]
}

// Keys returns all the colors in the set in unspecified order
func (s *CountedSet) Keys() []rgb {
	keys := make([]rgb, 0, len(s.m))

	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

// SortedSet returns the entries (Color, Count) ordered from greatest count to least
func (s *CountedSet) SortedSet() []CountedEntry {
	list := make([]CountedEntry, 0, len(s.m))

	for color, cnt := range s.m {
		list = append(list, CountedEntry{color, cnt})
	}

	sort.Sort(ByCount(list))
	return list
}
