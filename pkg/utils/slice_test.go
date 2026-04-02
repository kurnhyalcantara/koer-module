package utils

import (
	"testing"
)

func TestContains(t *testing.T) {
	if !Contains([]int{1, 2, 3}, 2) {
		t.Error("expected true")
	}
	if Contains([]int{1, 2, 3}, 4) {
		t.Error("expected false")
	}
}

func TestMap(t *testing.T) {
	result := Map([]int{1, 2, 3}, func(x int) int { return x * 2 })
	if len(result) != 3 || result[0] != 2 || result[1] != 4 || result[2] != 6 {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestFilter(t *testing.T) {
	result := Filter([]int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 })
	if len(result) != 2 || result[0] != 2 || result[1] != 4 {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestUnique(t *testing.T) {
	result := Unique([]int{1, 2, 2, 3, 1})
	if len(result) != 3 {
		t.Errorf("expected 3 unique elements, got %v", result)
	}
}
