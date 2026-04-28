package differ

import (
	"testing"
)

func buildDeltaTargets() map[string][]Result {
	return map[string][]Result{
		"staging": {
			{Key: "DB_HOST", Status: StatusAdded},
			{Key: "API_KEY", Status: StatusChanged},
			{Key: "OLD_VAR", Status: StatusRemoved},
		},
		"prod": {
			{Key: "DB_HOST", Status: StatusAdded},
			{Key: "OLD_VAR", Status: StatusRemoved},
			{Key: "NEW_FLAG", Status: StatusAdded},
		},
	}
}

func TestDelta_CountsChanges(t *testing.T) {
	entries := Delta(buildDeltaTargets(), DefaultDeltaOptions())
	found := map[string]DeltaEntry{}
	for _, e := range entries {
		found[e.Key] = e
	}

	if found["DB_HOST"].Adds != 2 {
		t.Errorf("DB_HOST adds: want 2, got %d", found["DB_HOST"].Adds)
	}
	if found["OLD_VAR"].Removals != 2 {
		t.Errorf("OLD_VAR removals: want 2, got %d", found["OLD_VAR"].Removals)
	}
	if found["API_KEY"].Changes != 1 {
		t.Errorf("API_KEY changes: want 1, got %d", found["API_KEY"].Changes)
	}
}

func TestDelta_NetCalculation(t *testing.T) {
	entries := Delta(buildDeltaTargets(), DefaultDeltaOptions())
	found := map[string]DeltaEntry{}
	for _, e := range entries {
		found[e.Key] = e
	}
	if found["DB_HOST"].Net != 2 {
		t.Errorf("DB_HOST net: want 2, got %d", found["DB_HOST"].Net)
	}
	if found["OLD_VAR"].Net != -2 {
		t.Errorf("OLD_VAR net: want -2, got %d", found["OLD_VAR"].Net)
	}
}

func TestDelta_MinNetFilters(t *testing.T) {
	opts := DefaultDeltaOptions()
	opts.MinNet = 2
	entries := Delta(buildDeltaTargets(), opts)
	for _, e := range entries {
		net := e.Net
		if net < 0 {
			net = -net
		}
		if net < 2 {
			t.Errorf("entry %q net %d below MinNet 2", e.Key, e.Net)
		}
	}
}

func TestDelta_Empty(t *testing.T) {
	entries := Delta(map[string][]Result{}, DefaultDeltaOptions())
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d entries", len(entries))
	}
}

func TestDelta_IncludeZero(t *testing.T) {
	targets := map[string][]Result{
		"a": {{Key: "X", Status: StatusUnchanged}},
	}
	opts := DefaultDeltaOptions()
	opts.IncludeZero = true
	entries := Delta(targets, opts)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry with IncludeZero, got %d", len(entries))
	}
}

func TestHasDeltaEntries(t *testing.T) {
	if HasDeltaEntries(nil) {
		t.Error("expected false for nil")
	}
	if !HasDeltaEntries([]DeltaEntry{{Key: "X"}}) {
		t.Error("expected true for non-empty")
	}
}
