package differ

import (
	"testing"
)

func TestSquash_Empty(t *testing.T) {
	out := Squash()
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestSquash_SingleSet(t *testing.T) {
	set := []Result{
		{Key: "A", Status: "added", New: "1"},
		{Key: "B", Status: "unchanged", Old: "x", New: "x"},
	}
	out := Squash(set)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestSquash_ChangedWinsOverAdded(t *testing.T) {
	set1 := []Result{{Key: "A", Status: "added", New: "1"}}
	set2 := []Result{{Key: "A", Status: "changed", Old: "1", New: "2"}}
	out := Squash(set1, set2)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Status != "changed" {
		t.Errorf("expected changed, got %s", out[0].Status)
	}
}

func TestSquash_UnchangedDoesNotOverwrite(t *testing.T) {
	set1 := []Result{{Key: "X", Status: "removed", Old: "v"}}
	set2 := []Result{{Key: "X", Status: "unchanged", Old: "v", New: "v"}}
	out := Squash(set1, set2)
	if out[0].Status != "removed" {
		t.Errorf("expected removed, got %s", out[0].Status)
	}
}

func TestSquash_SortedOutput(t *testing.T) {
	set := []Result{
		{Key: "Z", Status: "added"},
		{Key: "A", Status: "added"},
		{Key: "M", Status: "unchanged"},
	}
	out := Squash(set)
	if out[0].Key != "A" || out[1].Key != "M" || out[2].Key != "Z" {
		t.Errorf("unexpected order: %v", out)
	}
}

func TestHasSquashedChanges_True(t *testing.T) {
	results := []SquashResult{{Key: "A", Status: "added"}}
	if !HasSquashedChanges(results) {
		t.Error("expected true")
	}
}

func TestHasSquashedChanges_False(t *testing.T) {
	results := []SquashResult{{Key: "A", Status: "unchanged"}}
	if HasSquashedChanges(results) {
		t.Error("expected false")
	}
}
