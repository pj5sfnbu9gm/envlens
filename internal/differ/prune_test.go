package differ

import (
	"testing"
)

func buildPruneResults() []Result {
	return []Result{
		{Key: "APP_HOST", Status: StatusUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "APP_PORT", Status: StatusChanged, OldValue: "8080", NewValue: "9090"},
		{Key: "DB_URL", Status: StatusAdded, OldValue: "", NewValue: "postgres://db"},
		{Key: "LEGACY_KEY", Status: StatusRemoved, OldValue: "old", NewValue: ""},
		{Key: "SECRET_TOKEN", Status: StatusChanged, OldValue: "abc", NewValue: "xyz"},
	}
}

func TestPrune_DefaultRemovesUnchanged(t *testing.T) {
	results := buildPruneResults()
	out := Prune(results, DefaultPruneOptions())
	for _, r := range out {
		if r.Status == StatusUnchanged {
			t.Errorf("expected unchanged entry %q to be pruned", r.Key)
		}
	}
	if len(out) != 4 {
		t.Errorf("expected 4 results, got %d", len(out))
	}
}

func TestPrune_RemoveAdded(t *testing.T) {
	results := buildPruneResults()
	opts := PruneOptions{RemoveAdded: true}
	out := Prune(results, opts)
	for _, r := range out {
		if r.Status == StatusAdded {
			t.Errorf("expected added entry %q to be pruned", r.Key)
		}
	}
}

func TestPrune_RemoveByExactKey(t *testing.T) {
	results := buildPruneResults()
	opts := PruneOptions{Keys: []string{"DB_URL", "LEGACY_KEY"}}
	out := Prune(results, opts)
	for _, r := range out {
		if r.Key == "DB_URL" || r.Key == "LEGACY_KEY" {
			t.Errorf("expected key %q to be pruned", r.Key)
		}
	}
	if len(out) != 3 {
		t.Errorf("expected 3 results, got %d", len(out))
	}
}

func TestPrune_RemoveByPrefix(t *testing.T) {
	results := buildPruneResults()
	opts := PruneOptions{Prefixes: []string{"SECRET_"}}
	out := Prune(results, opts)
	for _, r := range out {
		if r.Key == "SECRET_TOKEN" {
			t.Errorf("expected SECRET_TOKEN to be pruned by prefix")
		}
	}
	if len(out) != 4 {
		t.Errorf("expected 4 results, got %d", len(out))
	}
}

func TestPrune_EmptyResults(t *testing.T) {
	out := Prune(nil, DefaultPruneOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d", len(out))
	}
}

func TestHasPruneResults_True(t *testing.T) {
	if !HasPruneResults([]Result{{Key: "X"}}) {
		t.Error("expected true")
	}
}

func TestHasPruneResults_False(t *testing.T) {
	if HasPruneResults(nil) {
		t.Error("expected false")
	}
}
