package zset

import (
	"reflect"
	"testing"
)

var members = []string{"hello", "secctan", "how", "are", "you"}

func createTestTree() *RBTree {
	tree := NewRBTree()
	for _, member := range members {
		tree.insert(member)
	}
	return tree
}

func TestInsert(t *testing.T) {
	tree := createTestTree()

	got := tree.members()
	want := []string{"are", "hello", "how", "secctan", "you"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestDelete(t *testing.T) {
	tree := createTestTree()
	tree.delete("hello")

	got := tree.members()
	want := []string{"are", "how", "secctan", "you"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestSearch(t *testing.T) {
	tree := createTestTree()

	_, got := tree.search("secctan")
	want := true
	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

	_, got = tree.search("world")
	want = false
	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}
}
