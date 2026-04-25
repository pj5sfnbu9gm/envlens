package differ

import (
	"testing"
)

func buildTrendWindows() []map[string][]DiffResult {
	// Three windows; DB_URL changes in all three, API_KEY only in window 0.
	w0 := map[string][]DiffResult{
		"prod": {
			{Key: "DB_URL", Status: StatusChanged, OldValue: "a", NewValue: "b"},
			{Key: "API_KEY", Status: StatusChanged, OldValue: "x", NewValue: "y"},
		},
	}
	w1 := map[string][]DiffResult{
		"prod": {
			{Key: "DB_URL", Status: StatusChanged, OldValue: "b", NewValue: "c"},
			{Key: "API_KEY", Status: StatusUnchanged, OldValue: "y", NewValue: "y"},
		},
	}
	w2 := map[string][]DiffResult{
		"prod": {
			{Key: "DB_URL", Status: StatusChanged, OldValue: "c", NewValue: "d"},
			{Key: "API_KEY", Status: StatusUnchanged, OldValue: "y", NewValue: "y"},
		},
	}
	return []map[string][]DiffResult{w0, w1, w2}
}

func TestTrend_Empty(t *testing.T) {
	result := Trend(nil, DefaultTrendOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d entries", len(result))
	}
}

func TestTrend_TotalCounts(t *testing.T) {
	windows := buildTrendWindows()
	entries := Trend(windows, DefaultTrendOptions())

	var dbEntry *TrendEntry
	for i := range entries {
		if entries[i].Key == "DB_URL" {
			dbEntry = &entries[i]
		}
	}
	if dbEntry == nil {
		t.Fatal("expected DB_URL entry")
	}
	if dbEntry.Total != 3 {
		t.Errorf("expected total 3, got %d", dbEntry.Total)
	}
}

func TestTrend_MinTotalFilters(t *testing.T) {
	windows := buildTrendWindows()
	opts := DefaultTrendOptions()
	opts.MinTotal = 2
	entries := Trend(windows, opts)
	for _, e := range entries {
		if e.Total < 2 {
			t.Errorf("entry %q has total %d below min", e.Key, e.Total)
		}
	}
}

func TestTrend_DirectionUp(t *testing.T) {
	// Craft windows where changes increase over time.
	w0 := map[string][]DiffResult{"t": {{Key: "X", Status: StatusUnchanged}}}
	w1 := map[string][]DiffResult{"t": {{Key: "X", Status: StatusChanged, OldValue: "a", NewValue: "b"}}}
	w2 := map[string][]DiffResult{"t": {{Key: "X", Status: StatusChanged, OldValue: "b", NewValue: "c"}}}
	windows := []map[string][]DiffResult{w0, w1, w2}

	opts := DefaultTrendOptions()
	entries := Trend(windows, opts)
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if entries[0].Direction != TrendUp {
		t.Errorf("expected TrendUp, got %s", entries[0].Direction)
	}
}

func TestTrend_IncludeStable(t *testing.T) {
	w0 := map[string][]DiffResult{"t": {{Key: "STABLE", Status: StatusUnchanged}}}
	windows := []map[string][]DiffResult{w0}

	opts := DefaultTrendOptions()
	opts.IncludeStable = true
	opts.MinTotal = 0
	entries := Trend(windows, opts)
	if len(entries) == 0 {
		t.Fatal("expected stable entry to be included")
	}
}

func TestHasTrendChanges_True(t *testing.T) {
	entries := []TrendEntry{{Key: "K", Total: 1}}
	if !HasTrendChanges(entries) {
		t.Error("expected true")
	}
}

func TestHasTrendChanges_False(t *testing.T) {
	entries := []TrendEntry{{Key: "K", Total: 0}}
	if HasTrendChanges(entries) {
		t.Error("expected false")
	}
}
