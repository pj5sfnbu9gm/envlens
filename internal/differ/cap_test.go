package differ

import (
	"testing"
)

func buildCapResults() map[string][]Result {
	return map[string][]Result{
		"staging": {
			{Key: "A", Status: StatusUnchanged},
			{Key: "B", Status: StatusAdded},
			{Key: "C", Status: StatusChanged, OldValue: "x", NewValue: "y"},
			{Key: "D", Status: StatusRemoved},
			{Key: "E", Status: StatusAdded},
		},
	}
}

func TestCap_NoLimit(t *testing.T) {
	results := buildCapResults()
	opts := DefaultCapOptions() // MaxPerTarget == 0
	out := Cap(results, opts)

	// Unchanged filtered, 4 entries remain
	if got := len(out["staging"]); got != 4 {
		t.Fatalf("expected 4 entries, got %d", got)
	}
}

func TestCap_LimitsToMax(t *testing.T) {
	results := buildCapResults()
	opts := DefaultCapOptions()
	opts.MaxPerTarget = 2
	out := Cap(results, opts)

	if got := len(out["staging"]); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}

func TestCap_PriorityOrder(t *testing.T) {
	results := buildCapResults()
	opts := DefaultCapOptions()
	opts.MaxPerTarget = 3
	out := Cap(results, opts)

	entries := out["staging"]
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	// First must be the Changed entry
	if entries[0].Status != StatusChanged {
		t.Errorf("expected first entry to be Changed, got %s", entries[0].Status)
	}
	// Next two must be Added
	for _, e := range entries[1:] {
		if e.Status != StatusAdded {
			t.Errorf("expected Added, got %s", e.Status)
		}
	}
}

func TestCap_IncludeUnchanged(t *testing.T) {
	results := buildCapResults()
	opts := CapOptions{MaxPerTarget: 10, IgnoreUnchanged: false}
	out := Cap(results, opts)

	if got := len(out["staging"]); got != 5 {
		t.Fatalf("expected 5 entries (all), got %d", got)
	}
}

func TestCap_EmptyResults(t *testing.T) {
	out := Cap(map[string][]Result{}, DefaultCapOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output")
	}
}

func TestHasCapResults_True(t *testing.T) {
	results := map[string][]Result{"prod": {{Key: "X", Status: StatusAdded}}}
	if !HasCapResults(results) {
		t.Error("expected true")
	}
}

func TestHasCapResults_False(t *testing.T) {
	results := map[string][]Result{"prod": {}}
	if HasCapResults(results) {
		t.Error("expected false")
	}
}
