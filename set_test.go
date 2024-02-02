package main

import "testing"

var members = []string{"hello", "world"}

func createTestSet() Set {
	set := NewSet()
	for _, member := range members {
		set.Add(member)
	}
	return set
}

func TestSize(t *testing.T) {
	set := createTestSet()

	got := set.Size()
	want := len(members)
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestExists(t *testing.T) {
	set := createTestSet()

	got := set.Exists("hello")
	want := true
	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

	got = set.Exists("secctan")
	want = false
	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}
}