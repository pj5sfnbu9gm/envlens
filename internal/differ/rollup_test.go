package differ

import (
	"testing"
)

func buildRollupTargets() map[string][]Result {
	return map[string][]Result{
		"staging": {
			{Key: "DB_HOST", Status: StatusChanged, OldValue: "a", NewValue: "b"},
			{Key: "API_KEY", Status: StatusAdded, NewValue: "xyz"},
			{Key: "PORT", Status: StatusUnchanged, OldValue: "8080", NewValue: "8080"},
		},
		"production": {
			{Key: "DB_HOST", Status: StatusChanged, OldValue: "a", NewValue: "c"},
			{Key: "PORT", Status: StatusRemoved, OldValue: "8080"},
		},
	}
}

func TestRollup_CountsPerKey(t *testing.T) {
	entries := Rollup(buildRollupTargets(), DefaultRollupOptions())
	found := map[string]RollupEntry{}
	for _, e := range entries {
		found[e.Key] = e
	}
	if found["DB_HOST"].Changed != 2 {
		t.Errorf("expected DB_HOST Changed=2, got %d", found["DB_HOST"].Changed)
	}
	if found["API_KEY"].Added != 1 {
		t.Errorf("expected API_KEY Added=1, got %d", found["API_KEY"].Added)
	}
	if found["PORT"].Removed != 1 {
		t.Errorf("expected PORT Removed=1, got %d", found["PORT"].Removed)
	}
}

func TestRollup_FiltersUnchangedByDefault(t *testing.T) {
	entries := Rollup(buildRollupTargets(), DefaultRollupOptions())
	for _, e := range entries {
		if e.Key == "PORT" && e.Total == 0 {
			t.Error("PORT should appear because it has a removal")
		}
	}
}

func TestRollup_MinChangesZero_IncludesUnchanged(t *testing.T) {
	opts := RollupOptions{MinChanges: 0}
	entries := Rollup(buildRollupTargets(), opts)
	keys := map[string]bool{}
	for _, e := range entries {
		keys[e.Key] = true
	}
	if !keys["PORT"] {
		t.Error("expected PORT with MinChanges=0")
	}
}

func TestRollup_SortedByTotalDesc(t *testing.T) {
	entries := Rollup(buildRollupTargets(), DefaultRollupOptions())
	for i := 1; i < len(entries); i++ {
		if entries[i].Total > entries[i-1].Total {
			t.Errorf("entries not sorted desc at index %d", i)
		}
	}
}

func TestRollup_Empty(t *testing.T) {
	entries := Rollup(map[string][]Result{}, DefaultRollupOptions())
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestHasRollupChanges_True(t *testing.T) {
	entries := Rollup(buildRollupTargets(), DefaultRollupOptions())
	if !HasRollupChanges(entries) {
		t.Error("expected HasRollupChanges=true")
	}
}

func TestHasRollupChanges_False(t *testing.T) {
	if HasRollupChanges([]RollupEntry{}) {
		t.Error("expected HasRollupChanges=false for empty slice")
	}
}
