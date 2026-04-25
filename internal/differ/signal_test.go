package differ

import (
	"testing"
)

func buildSignalTargets() map[string][]Result {
	return map[string][]Result{
		"prod": {
			{Key: "DB_HOST", Status: StatusChanged, OldValue: "a", NewValue: "b"},
			{Key: "API_KEY", Status: StatusChanged, OldValue: "x", NewValue: "y"},
			{Key: "LOG_LEVEL", Status: StatusUnchanged, OldValue: "info", NewValue: "info"},
		},
		"staging": {
			{Key: "DB_HOST", Status: StatusChanged, OldValue: "a", NewValue: "c"},
			{Key: "TIMEOUT", Status: StatusAdded, OldValue: "", NewValue: "30"},
		},
		"dev": {
			{Key: "DB_HOST", Status: StatusChanged, OldValue: "a", NewValue: "d"},
			{Key: "API_KEY", Status: StatusRemoved, OldValue: "z", NewValue: ""},
		},
	}
}

func TestSignal_Empty(t *testing.T) {
	result := Signal(map[string][]Result{}, DefaultSignalOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

func TestSignal_CountsAcrossTargets(t *testing.T) {
	entries := Signal(buildSignalTargets(), DefaultSignalOptions())
	if len(entries) == 0 {
		t.Fatal("expected entries")
	}
	// DB_HOST changed in 3 targets
	if entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST first, got %s", entries[0].Key)
	}
	if entries[0].ChangeCount != 3 {
		t.Errorf("expected count 3, got %d", entries[0].ChangeCount)
	}
}

func TestSignal_MinTargetsFilters(t *testing.T) {
	opts := DefaultSignalOptions()
	opts.MinTargets = 3
	entries := Signal(buildSignalTargets(), opts)
	for _, e := range entries {
		if e.ChangeCount < 3 {
			t.Errorf("entry %s has count %d below min", e.Key, e.ChangeCount)
		}
	}
}

func TestSignal_TopN(t *testing.T) {
	opts := DefaultSignalOptions()
	opts.MinTargets = 1
	opts.TopN = 1
	entries := Signal(buildSignalTargets(), opts)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestSignal_TargetsSorted(t *testing.T) {
	entries := Signal(buildSignalTargets(), DefaultSignalOptions())
	for _, e := range entries {
		for i := 1; i < len(e.Targets); i++ {
			if e.Targets[i] < e.Targets[i-1] {
				t.Errorf("targets not sorted for key %s", e.Key)
			}
		}
	}
}

func TestHasSignals_True(t *testing.T) {
	entries := Signal(buildSignalTargets(), DefaultSignalOptions())
	if !HasSignals(entries) {
		t.Error("expected HasSignals true")
	}
}

func TestHasSignals_False(t *testing.T) {
	if HasSignals(nil) {
		t.Error("expected HasSignals false")
	}
}
