package differ

import (
	"testing"
)

func TestBloom_Empty(t *testing.T) {
	results := Bloom(nil, DefaultBloomOptions())
	if len(results) != 0 {
		t.Fatalf("expected empty, got %d entries", len(results))
	}
}

func TestBloom_AllPresent(t *testing.T) {
	targets := map[string]map[string]string{
		"dev":  {"FOO": "1", "BAR": "2"},
		"prod": {"FOO": "1", "BAR": "2"},
	}
	results := Bloom(targets, DefaultBloomOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
	for _, e := range results {
		if len(e.AbsentIn) != 0 {
			t.Errorf("key %q should be present everywhere", e.Key)
		}
	}
}

func TestBloom_MissingInOne(t *testing.T) {
	targets := map[string]map[string]string{
		"dev":  {"FOO": "1", "BAR": "2"},
		"prod": {"FOO": "1"},
	}
	results := Bloom(targets, DefaultBloomOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
	var barEntry *BloomEntry
	for i := range results {
		if results[i].Key == "BAR" {
			barEntry = &results[i]
		}
	}
	if barEntry == nil {
		t.Fatal("expected BAR entry")
	}
	if len(barEntry.AbsentIn) != 1 || barEntry.AbsentIn[0] != "prod" {
		t.Errorf("expected BAR absent in prod, got %v", barEntry.AbsentIn)
	}
}

func TestBloom_MinPresenceFilter(t *testing.T) {
	targets := map[string]map[string]string{
		"dev":     {"FOO": "1"},
		"staging": {"FOO": "1"},
		"prod":    {"BAR": "2"},
	}
	opts := BloomOptions{MinPresence: 2}
	results := Bloom(targets, opts)
	if len(results) != 1 || results[0].Key != "FOO" {
		t.Fatalf("expected only FOO (present in 2 targets), got %v", results)
	}
}

func TestBloom_SortedOutput(t *testing.T) {
	targets := map[string]map[string]string{
		"a": {"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"},
	}
	results := Bloom(targets, DefaultBloomOptions())
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, e := range results {
		if e.Key != expected[i] {
			t.Errorf("index %d: want %s got %s", i, expected[i], e.Key)
		}
	}
}

func TestHasBloomGaps_True(t *testing.T) {
	entries := []BloomEntry{
		{Key: "FOO", PresentIn: []string{"dev"}, AbsentIn: []string{"prod"}},
	}
	if !HasBloomGaps(entries) {
		t.Error("expected gaps")
	}
}

func TestHasBloomGaps_False(t *testing.T) {
	entries := []BloomEntry{
		{Key: "FOO", PresentIn: []string{"dev", "prod"}, AbsentIn: nil},
	}
	if HasBloomGaps(entries) {
		t.Error("expected no gaps")
	}
}
