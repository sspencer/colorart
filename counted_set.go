package colorart

import (
	"fmt"
	"sort"
)

// CountedEntry is for use by sorting class
type CountedEntry struct {
	Name  string
	Count int
}

// ByCount is the type used to sort
type ByCount []CountedEntry

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].Count > a[j].Count }

func (e CountedEntry) String() string {
	return fmt.Sprintf("%s: %d", e.Name, e.Count)
}

// CountedSet counts the number of times each object (string) is added to the set.
// The set is not thread safe.
type CountedSet struct {
	m map[string]int
}

//---------------------------

// NewCountedSet creates a new CountedSet of the specified size.
func NewCountedSet(size int) *CountedSet {
	s := &CountedSet{}
	s.m = make(map[string]int, size)
	return s
}

// Add adds an object to the set.
func (s *CountedSet) Add(name string) {
	s.m[name]++
}

func (s *CountedSet) AddCount(name string, count int) {
	s.m[name] = count
}

// Size returns the number of objects in the set.
func (s *CountedSet) Size() int {
	return len(s.m)
}

// Count returns the number of times the specified object has been added to the set.
func (s *CountedSet) Count(name string) int {
	return s.m[name]
}

// Remove decrements the number of times the specified object has been added to the set.
func (s *CountedSet) Remove(name string) {
	count, ok := s.m[name]
	if ok {
		if count > 1 {
			s.m[name]--
		} else {
			delete(s.m, name)
		}
	}
}

// RemoveAll removes the specified object completely from the set (Count goes to 0)
func (s *CountedSet) RemoveAll(name string) {
	delete(s.m, name)
}

// Keys returns all the names in the set in unspecified order
func (s *CountedSet) Keys() []string {
	var keys []string
	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

// SortedSet returns the entries (Name, Count) order from greatest count to least
func (s *CountedSet) SortedSet() []CountedEntry {
	list := make([]CountedEntry, 0, len(s.m))

	for name, cnt := range s.m {
		list = append(list, CountedEntry{name, cnt})
	}

	sort.Sort(ByCount(list))
	return list
}
