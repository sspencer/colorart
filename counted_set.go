package colorart

import (
	"fmt"
	"sort"
)

const maxComponent = 255

// rgb is RGB components, 0-255, used as the map key
type rgb [3]byte

// countedSet counts the number of times each object (string) is added to the set.
// The set is not thread safe.
type countedSet struct {
	set map[rgb]int

	// detune colors so colors within a few shades of each other
	// map to the same color.  Makes algorthim much faster.
	// 0 is no change, original color
	// 1 divides by 2, multiplies by 2 (so 02, 03 map to same color)
	// 2 divides by 4, multiplies by 4 (so 04, 05, 06, 07 map to same)
	// 3 divides by 8, multiplies by 8
	// Don't go much beyond 3...
	shift uint8
}

// countedEntry is for use by sorting class
type countedEntry struct {
	color rgb
	count int
}

// byCount is the type used to sort
type byCount []countedEntry

func (a byCount) Len() int           { return len(a) }
func (a byCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCount) Less(i, j int) bool { return a[i].count > a[j].count }

func (e countedEntry) String() string {
	return fmt.Sprintf("#%02x%02x%02x:%d", e.color[0], e.color[1], e.color[2], e.count)
}

//---------------------------

// newCountedSet creates a new countedSet of the specified size.
func newCountedSet(size, colorShift int) countedSet {
	set := make(map[rgb]int, size)
	shift := uint8(colorShift)
	return countedSet{set, shift}
}

// addPixel converts pixel to [3]byte rgb and counts unique colors
func (s countedSet) addPixel(p pixel) {

	ri := uint8(maxComponent*p.R) >> s.shift << s.shift
	gi := uint8(maxComponent*p.G) >> s.shift << s.shift
	bi := uint8(maxComponent*p.B) >> s.shift << s.shift

	color := rgb{ri, gi, bi}

	s.set[color]++
}

// merge other counted set into this one.
func (s countedSet) merge(o countedSet) {
	for color, cnt := range o.set {
		s.set[color] += cnt
	}
}

// addCount adds color into set
func (s countedSet) addCount(color rgb, count int) {
	s.set[color] += count
}

// count returns the number of times the specified object has been added to the set.
func (s countedSet) count(color rgb) int {
	return s.set[color]
}

// sortedSet returns the entries (Color, Count) ordered from greatest count to least
func (s countedSet) sortedSet() []countedEntry {

	list := make([]countedEntry, len(s.set))
	it := 0
	for color, cnt := range s.set {
		list[it] = countedEntry{color, cnt}
		it++
	}

	sort.Sort(byCount(list))
	return list
}
