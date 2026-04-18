package differ

import (
	"testing"
)

func buildMultiResults() map[string][]Result {
	return map[string][]Result{
		"staging": {
			{Key: "FOO", Status: StatusChanged},
			{Key: "BAR", Status: StatusChanged},
			{Key: "BAZ", Status: StatusUnchanged},
		},
		"prod": {
			{Key: "FOO", Status: StatusChanged},
			{Key: "BAZ", Status: StatusAdded},
		},
	}
}

func TestRank_CountsChanges(t *testing.T) {
	results := buildMultiResults()
	entries := Rank(results, DefaultRankOptions())
	if len(entries) == 0 {
		t.Fatal("expected entries")
	}
	if entries[0].Key != "FOO" || entries[0].Changes != 2 {
		t.Errorf("expected FOO with 2 changes, got %+v", entries[0])
	}
}

func TestRank_MinChangesFilter(t *testing.T) {
	results := buildMultiResults()
	opts := DefaultRankOptions()
	opts.MinChanges = 2
	entries := Rank(results, opts)
	for _, e := range entries {
		if e.Changes < 2 {
			t.Errorf("entry %s has %d changes, below min", e.Key, e.Changes)
		}
	}
}

func TestRank_TopN(t *testing.T) {
	results := buildMultiResults()
	opts := DefaultRankOptions()
	opts.TopN = 1
	entries := Rank(results, opts)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestRank_Empty(t *testing.T) {
	entries := Rank(map[string][]Result{}, DefaultRankOptions())
	if HasRankResults(entries) {
		t.Error("expected no results")
	}
}

func TestRank_SortedByKeyWhenTied(t *testing.T) {
	results := map[string][]Result{
		"t1": {
			{Key: "ZEBRA", Status: StatusChanged},
			{Key: "ALPHA", Status: StatusChanged},
		},
	}
	entries := Rank(results, DefaultRankOptions())
	if entries[0].Key != "ALPHA" {
		t.Errorf("expected ALPHA first, got %s", entries[0].Key)
	}
}
