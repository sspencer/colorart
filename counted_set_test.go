package colorart

import "testing"

func TestSet(t *testing.T) {
	s := NewCountedSet(10)
	s.Add("One")
	s.Add("Two")
	s.Add("Two")

	if s.Count("One") != 1 {
		t.Error("Add: incorrect One count")
	}
	if s.Count("Two") != 2 {
		t.Error("Add: incorrect Two count")
	}
	if s.Count("Zero") != 0 {
		t.Error("Add: incorrect Zero count")
	}
	if s.Size() != 2 {
		// Only "One" and "Two" should be in Set
		t.Error("Size: incorrect Size after Add")
	}

	var hasOne, hasTwo bool
	for _, k := range s.Keys() {
		if k == "One" {
			hasOne = true
		}
		if k == "Two" {
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

	if sorted[0].Name != "Two" && sorted[1].Name != "One" {
		t.Error("Sorted set is not sorted")
	}

	s.Remove("Two")
	if s.Count("Two") != 1 {
		t.Error("Remove: incorrect Count after Remove")
	}

	s.RemoveAll("Two")
	if s.Count("Two") != 0 {
		t.Error("Remove: incorrect Count after RemoveAll")
	}
	if s.Size() != 1 {
		// Only "One" and "Two" should be in Set
		t.Error("Size: incorrect Size after RemoveAll")
	}

	s.Remove("SomeRandomObject")
	if s.Size() != 1 {
		t.Error("Size: incorrect Size after Remove of unknown object")
	}

	s.RemoveAll("SomeOtherRandomObject")
	if s.Size() != 1 {
		t.Error("Size: incorrect Size after Remove of unknown object")
	}

}
