package differ

import (
	"testing"
)

func buildCensusTargets() map[string]map[string]string {
	return map[string]map[string]string{
		"prod": {"DB_HOST": "db.prod", "API_KEY": "secret", "TIMEOUT": "30"},
		"staging": {"DB_HOST": "db.staging", "API_KEY": "s3cr3t"},
		"dev": {"DB_HOST": "localhost", "DEBUG": "true"},
	}
}

func TestCensus_Empty(t *testing.T) {
	result := Census(nil, DefaultCensusOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestCensus_CountsTargets(t *testing.T) {
	targets := buildCensusTargets()
	out := Census(targets, DefaultCensusOptions())

	find := func(key string) *CensusEntry {
		for i := range out {
			if out[i].Key == key {
				return &out[i]
			}
		}
		return nil
	}

	dbHost := find("DB_HOST")
	if dbHost == nil {
		t.Fatal("DB_HOST not found")
	}
	if dbHost.Count != 3 {
		t.Errorf("expected count 3, got %d", dbHost.Count)
	}

	debug := find("DEBUG")
	if debug == nil {
		t.Fatal("DEBUG not found")
	}
	if debug.Count != 1 {
		t.Errorf("expected count 1, got %d", debug.Count)
	}
}

func TestCensus_CoverageCalculation(t *testing.T) {
	targets := buildCensusTargets()
	out := Census(targets, DefaultCensusOptions())
	for _, e := range out {
		expected := float64(e.Count) / float64(len(targets))
		if e.Coverage != expected {
			t.Errorf("key %s: coverage %f, want %f", e.Key, e.Coverage, expected)
		}
	}
}

func TestCensus_MinCoverageFilter(t *testing.T) {
	targets := buildCensusTargets()
	opts := DefaultCensusOptions()
	opts.MinCoverage = 0.9 // only keys in all 3 targets
	out := Census(targets, opts)
	for _, e := range out {
		if e.Count < 3 {
			t.Errorf("key %s should have been filtered (count=%d)", e.Key, e.Count)
		}
	}
}

func TestCensus_ExcludeUniversal(t *testing.T) {
	targets := buildCensusTargets()
	opts := DefaultCensusOptions()
	opts.ExcludeUniversal = true
	out := Census(targets, opts)
	for _, e := range out {
		if e.Count == len(targets) {
			t.Errorf("universal key %s should be excluded", e.Key)
		}
	}
}

func TestCensus_SortedByCountDesc(t *testing.T) {
	targets := buildCensusTargets()
	out := Census(targets, DefaultCensusOptions())
	for i := 1; i < len(out); i++ {
		if out[i].Count > out[i-1].Count {
			t.Errorf("not sorted desc at index %d", i)
		}
	}
}

func TestHasCensusGaps_True(t *testing.T) {
	entries := []CensusEntry{{Key: "X", Count: 2}}
	if !HasCensusGaps(entries, 3) {
		t.Error("expected gaps")
	}
}

func TestHasCensusGaps_False(t *testing.T) {
	entries := []CensusEntry{{Key: "X", Count: 3}}
	if HasCensusGaps(entries, 3) {
		t.Error("expected no gaps")
	}
}
