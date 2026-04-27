package differ

import (
	"testing"
)

func buildHeatmapResults() map[string][]Result {
	return map[string][]Result{
		"staging": {
			{Key: "DB_HOST", Status: StatusChanged, OldValue: "a", NewValue: "b"},
			{Key: "PORT", Status: StatusUnchanged, OldValue: "8080", NewValue: "8080"},
			{Key: "SECRET", Status: StatusAdded, NewValue: "x"},
		},
		"prod": {
			{Key: "DB_HOST", Status: StatusChanged, OldValue: "b", NewValue: "c"},
			{Key: "SECRET", Status: StatusRemoved, OldValue: "x"},
		},
	}
}

func TestHeatmap_CountsChanges(t *testing.T) {
	results := buildHeatmapResults()
	entries := Heatmap(results, DefaultHeatmapOptions())
	if len(entries) == 0 {
		t.Fatal("expected entries")
	}
	// DB_HOST changed in both targets → 2
	var dbEntry *HeatmapEntry
	for i := range entries {
		if entries[i].Key == "DB_HOST" {
			dbEntry = &entries[i]
		}
	}
	if dbEntry == nil {
		t.Fatal("DB_HOST not found")
	}
	if dbEntry.Changes != 2 {
		t.Errorf("expected 2 changes, got %d", dbEntry.Changes)
	}
}

func TestHeatmap_UnchangedExcluded(t *testing.T) {
	results := buildHeatmapResults()
	entries := Heatmap(results, DefaultHeatmapOptions())
	for _, e := range entries {
		if e.Key == "PORT" {
			t.Error("PORT should be excluded (unchanged)")
		}
	}
}

func TestHeatmap_MinChangesFilter(t *testing.T) {
	results := buildHeatmapResults()
	opts := DefaultHeatmapOptions()
	opts.MinChanges = 2
	entries := Heatmap(results, opts)
	for _, e := range entries {
		if e.Changes < 2 {
			t.Errorf("entry %q has only %d changes, should be filtered", e.Key, e.Changes)
		}
	}
}

func TestHeatmap_TopN(t *testing.T) {
	results := buildHeatmapResults()
	opts := DefaultHeatmapOptions()
	opts.TopN = 1
	entries := Heatmap(results, opts)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST as top entry, got %s", entries[0].Key)
	}
}

func TestHeatmap_Empty(t *testing.T) {
	entries := Heatmap(map[string][]Result{}, DefaultHeatmapOptions())
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d", len(entries))
	}
	if HasHeatmapEntries(entries) {
		t.Error("HasHeatmapEntries should be false")
	}
}

func TestHeatmap_TargetList(t *testing.T) {
	results := buildHeatmapResults()
	entries := Heatmap(results, DefaultHeatmapOptions())
	for _, e := range entries {
		if e.Key == "DB_HOST" {
			if len(e.Targets) != 2 {
				t.Errorf("expected 2 targets for DB_HOST, got %d", len(e.Targets))
			}
		}
	}
}
