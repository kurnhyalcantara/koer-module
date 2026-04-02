package utils

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	s := RandomString(16)
	if len(s) != 16 {
		t.Errorf("expected length 16, got %d", len(s))
	}
	s2 := RandomString(16)
	if s == s2 {
		t.Error("two random strings should not be equal")
	}
}

func TestRandomID(t *testing.T) {
	id := RandomID()
	if len(id) != 36 {
		t.Errorf("expected length 36, got %d: %s", len(id), id)
	}
	id2 := RandomID()
	if id == id2 {
		t.Error("two random IDs should not be equal")
	}
}
