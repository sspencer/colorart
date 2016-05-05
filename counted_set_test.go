package colorart

import "testing"

const defaultColorShift = 0

var (
	zero = rgb{0, 100, 200}
	one  = rgb{1, 11, 111}
	two  = rgb{2, 22, 222}
)

func TestSetAddPixel(t *testing.T) {
	s := newCountedSet(10, defaultColorShift)
	s.addPixel(pixel{0.2, 0.4, 0.6, 1})
	entries := s.sortedSet()
	e := entries[0]
	str := e.String()
	answer := "#336699:1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}

func TestSetCount(t *testing.T) {
	s := newCountedSet(10, defaultColorShift)

	s.set[one]++
	if s.count(one) != 1 {
		t.Error("Count should be 1")
	}

	s.set[one]++
	if s.count(one) != 2 {
		t.Error("Count should be 2")
	}

	s.set[one]++
	if s.count(one) != 3 {
		t.Error("Count should be 3")
	}

	if len(s.set) != 1 {
		// Only "one" and "two" should be in Set
		t.Error("len: incorrect size after Add")
	}
}

func TestSetMultiCount(t *testing.T) {
	s := newCountedSet(10, defaultColorShift)

	s.set[one]++
	s.set[two]++
	s.set[two]++

	if s.count(one) != 1 {
		t.Error("Add: incorrect 'one' count")
	}
	if s.count(two) != 2 {
		t.Error("Add: incorrect 'two' count")
	}
	if s.count(zero) != 0 {
		t.Error("Add: incorrect 'zero' count")
	}
	if len(s.set) != 2 {
		// Only "one" and "two" should be in Set
		t.Error("len: incorrect size after Add")
	}
}

func TestSetSortedSet(t *testing.T) {
	s := newCountedSet(10, defaultColorShift)

	s.set[one]++
	s.set[two]++
	s.set[two]++

	entries := s.sortedSet()
	if len(entries) != 2 {
		t.Error("Sorted set returned incorrect number of elements")
	}

	e := entries[0]
	if e.color != two && e.count != 2 {
		t.Error("First sorted entry (two) is incorrect")
	}

	e = entries[1]
	if e.color != one && e.count != 1 {
		t.Error("First sorted entry (two) is incorrect")
	}
}

func TestSetAddCount(t *testing.T) {
	s := newCountedSet(10, defaultColorShift)

	s.set[one]++
	if s.count(one) != 1 {
		t.Error("Incorrect count before AddCount")
	}

	// does not add new count to old ... replaces old with new
	s.addCount(one, 108)
	if s.count(one) != 109 {
		t.Error("Incorrect count after AddCount(1)")
	}

	s.addCount(two, 23456)
	if s.count(two) != 23456 {
		t.Error("Incorrect count after AddCount(2)")
	}
}

func TestSetString(t *testing.T) {
	s := newCountedSet(10, defaultColorShift)

	s.set[one]++
	entries := s.sortedSet()
	e := entries[0]
	str := e.String()

	answer := "#010b6f:1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}

func TestSetColorShift0(t *testing.T) {
	s := newCountedSet(10, 0)
	s.addPixel(pixel{0.2, 0.4, 0.6, 1})
	entries := s.sortedSet()
	e := entries[0]
	str := e.String()

	answer := "#336699:1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}

func TestSetColorShift1(t *testing.T) {
	s := newCountedSet(10, 1)
	s.addPixel(pixel{0.2, 0.4, 0.6, 1})
	entries := s.sortedSet()
	e := entries[0]
	str := e.String()

	answer := "#326698:1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}

func TestSetColorShift2(t *testing.T) {
	s := newCountedSet(10, 2)
	s.addPixel(pixel{0.2, 0.4, 0.6, 1})
	entries := s.sortedSet()
	e := entries[0]
	str := e.String()

	answer := "#306498:1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}

func TestSetColorShift3(t *testing.T) {
	s := newCountedSet(10, 3)
	s.addPixel(pixel{0.2, 0.4, 0.6, 1})
	entries := s.sortedSet()
	e := entries[0]
	str := e.String()

	answer := "#306098:1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}

func TestSetColorShift4(t *testing.T) {
	s := newCountedSet(10, 4)
	s.addPixel(pixel{0.2, 0.4, 0.6, 1})
	entries := s.sortedSet()
	e := entries[0]
	str := e.String()

	answer := "#306090:1"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}
