package differ

import (
	"testing"
)

func buildSliceResults() []Result {
	return []Result{
		{Key: "DB_HOST", Status: StatusChanged, OldValue: "old", NewValue: "new"},
		{Key: "DB_PORT", Status: StatusUnchanged, OldValue: "5432", NewValue: "5432"},
		{Key: "APP_ENV", Status: StatusAdded, NewValue: "prod"},
		{Key: "APP_DEBUG", Status: StatusRemoved, OldValue: "true"},
		{Key: "REDIS_URL", Status: StatusChanged, OldValue: "a", NewValue: "b"},
	}
}

func TestSlice_GroupsByPrefix(t *testing.T) {
	results := buildSliceResults()
	slices := Slice(results, DefaultSliceOptions())
	if len(slices) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(slices))
	}
	if slices[0].Prefix != "DB" {
		t.Errorf("expected first prefix DB, got %s", slices[0].Prefix)
	}
}

func TestSlice_IgnoreUnchanged(t *testing.T) {
	results := buildSliceResults()
	opts := DefaultSliceOptions()
	opts.IgnoreUnchanged = true
	slices := Slice(results, opts)
	for _, s := range slices {
		for _, r := range s.Results {
			if r.Status == StatusUnchanged {
				t.Errorf("unchanged result leaked into slice %s", s.Prefix)
			}
		}
	}
}

func TestSlice_IncludeUnchanged(t *testing.T) {
	results := buildSliceResults()
	opts := DefaultSliceOptions()
	opts.IgnoreUnchanged = false
	slices := Slice(results, opts)
	total := 0
	for _, s := range slices {
		total += len(s.Results)
	}
	if total != len(results) {
		t.Errorf("expected %d total results, got %d", len(results), total)
	}
}

func TestSlice_AllowedPrefixes(t *testing.T) {
	results := buildSliceResults()
	opts := DefaultSliceOptions()
	opts.Prefixes = []string{"DB", "APP"}
	opts.IgnoreUnchanged = false
	slices := Slice(results, opts)
	for _, s := range slices {
		if s.Prefix != "DB" && s.Prefix != "APP" && s.Prefix != "other" {
			t.Errorf("unexpected prefix %s", s.Prefix)
		}
	}
}

func TestSlice_Empty(t *testing.T) {
	slices := Slice(nil, DefaultSliceOptions())
	if slices != nil {
		t.Errorf("expected nil for empty input")
	}
}

func TestHasSliceChanges_True(t *testing.T) {
	results := buildSliceResults()
	slices := Slice(results, DefaultSliceOptions())
	if !HasSliceChanges(slices) {
		t.Error("expected HasSliceChanges to be true")
	}
}

func TestHasSliceChanges_False(t *testing.T) {
	if HasSliceChanges(nil) {
		t.Error("expected HasSliceChanges to be false for nil")
	}
}
