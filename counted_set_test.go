package colorart

import "testing"

var (
	zero = rgb{0, 100, 200}
	one  = rgb{1, 11, 111}
	two  = rgb{2, 22, 222}
)

func TestAddPixel(t *testing.T) {
	s := NewCountedSet(10)
	s.AddPixel(pixel{0.2, 0.4, 0.6, 1})
	entries := s.SortedSet()
	e := entries[0]
	str := e.String()
	answer := "336699: 1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}

func TestCount(t *testing.T) {
	s := NewCountedSet(10)

	s[one]++
	if s.Count(one) != 1 {
		t.Error("Count should be 1")
	}

	s[one]++
	if s.Count(one) != 2 {
		t.Error("Count should be 2")
	}

	s[one]++
	if s.Count(one) != 3 {
		t.Error("Count should be 3")
	}

	if len(s) != 1 {
		// Only "one" and "two" should be in Set
		t.Error("len: incorrect size after Add")
	}
}

func TestMultiCount(t *testing.T) {
	s := NewCountedSet(10)

	s[one]++
	s[two]++
	s[two]++

	if s.Count(one) != 1 {
		t.Error("Add: incorrect 'one' count")
	}
	if s.Count(two) != 2 {
		t.Error("Add: incorrect 'two' count")
	}
	if s.Count(zero) != 0 {
		t.Error("Add: incorrect 'zero' count")
	}
	if len(s) != 2 {
		// Only "one" and "two" should be in Set
		t.Error("len: incorrect size after Add")
	}
}

func TestSortedSet(t *testing.T) {
	s := NewCountedSet(10)

	s[one]++
	s[two]++
	s[two]++

	entries := s.SortedSet()
	if len(entries) != 2 {
		t.Error("Sorted set returned incorrect number of elements")
	}

	e := entries[0]
	if e.Color != two && e.Count != 2 {
		t.Error("First sorted entry (two) is incorrect")
	}

	e = entries[1]
	if e.Color != one && e.Count != 1 {
		t.Error("First sorted entry (two) is incorrect")
	}
}

func TestAddCount(t *testing.T) {
	s := NewCountedSet(10)

	s[one]++
	if s.Count(one) != 1 {
		t.Error("Incorrect count before AddCount")
	}

	// does not add new count to old ... replaces old with new
	s.AddCount(one, 108)
	if s.Count(one) != 109 {
		t.Error("Incorrect count after AddCount(1)")
	}

	s.AddCount(two, 23456)
	if s.Count(two) != 23456 {
		t.Error("Incorrect count after AddCount(2)")
	}
}

func TestString(t *testing.T) {
	s := NewCountedSet(10)

	s[one]++
	entries := s.SortedSet()
	e := entries[0]
	str := e.String()

	answer := "010b6f: 1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}
