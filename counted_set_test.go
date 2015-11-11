package colorart

import "testing"

func TestSet(t *testing.T) {
	s := NewCountedSet(10)
	one := rgb{51, 102, 153}
	two := rgb{153, 153, 153}
	zero := rgb{10, 20, 30}

	s.Add(one)
	s.Add(two)
	s.Add(two)

	if s.Count(one) != 1 {
		t.Error("Add: incorrect 'one' count")
	}
	if s.Count(two) != 2 {
		t.Error("Add: incorrect 'two' count")
	}
	if s.Count(zero) != 0 {
		t.Error("Add: incorrect 'zero' count")
	}
	if s.Size() != 2 {
		// Only "One" and "Two" should be in Set
		t.Error("Size: incorrect Size after Add")
	}

	var hasOne, hasTwo bool
	for _, k := range s.Keys() {
		if k == one {
			hasOne = true
		}
		if k == two {
			hasTwo = true
		}
	}

	if !hasOne {
		t.Error("Iteration does not have One")
	}

	if !hasTwo {
		t.Error("Iteration does not have Two")
	}

	sorted := s.SortedSet()
	if len(sorted) != 2 {
		t.Error("Sorted set does not have enough elements")
	}

	if sorted[0].Color != two && sorted[1].Color != one {
		t.Error("Sorted set is not sorted")
	}

	s.Remove(two)
	if s.Count(two) != 1 {
		t.Error("Remove: incorrect Count after Remove")
	}

	s.RemoveAll(two)
	if s.Count(two) != 0 {
		t.Error("Remove: incorrect Count after RemoveAll")
	}
	if s.Size() != 1 {
		// Only "One" and "Two" should be in Set
		t.Error("Size: incorrect Size after RemoveAll")
	}

	s.Remove(zero)
	if s.Size() != 1 {
		t.Error("Size: incorrect Size after Remove of unknown object")
	}

	s.RemoveAll(zero)
	if s.Size() != 1 {
		t.Error("Size: incorrect Size after Remove of unknown object")
	}
}
