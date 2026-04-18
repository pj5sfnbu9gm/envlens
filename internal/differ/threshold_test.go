package differ

import (
	"testing"
)

func makeTargetResults() map[string][]Result {
	return map[string][]Result{
		"prod": {
			{Key: "FOO", Status: StatusChanged, OldValue: "a", NewValue: "b"},
			{Key: "BAR", Status: StatusUnchanged, OldValue: "x", NewValue: "x"},
		},
		"staging": {
			{Key: "FOO", Status: StatusUnchanged, OldValue: "a", NewValue: "a"},
		},
		"dev": {
			{Key: "FOO", Status: StatusAdded, NewValue: "z"},
			{Key: "BAZ", Status: StatusRemoved, OldValue: "old"},
		},
	}
}

func TestApplyThreshold_DefaultFiltersUnchanged(t *testing.T) {
	results := ApplyThreshold(makeTargetResults(), DefaultThresholdOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(results))
	}
	for _, r := range results {
		if r.Target == "staging" {
			t.Errorf("staging should have been filtered out")
		}
	}
}

func TestApplyThreshold_MinChangesTwo(t *testing.T) {
	opts := DefaultThresholdOptions()
	opts.MinChanges = 2
	results := ApplyThreshold(makeTargetResults(), opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 target, got %d", len(results))
	}
	if results[0].Target != "dev" {
		t.Errorf("expected dev, got %s", results[0].Target)
	}
}

func TestApplyThreshold_IncludeUnchanged(t *testing.T) {
	opts := DefaultThresholdOptions()
	opts.IncludeUnchanged = true
	results := ApplyThreshold(makeTargetResults(), opts)
	for _, r := range results {
		if r.Target == "prod" && len(r.Results) != 2 {
			t.Errorf("expected 2 results for prod, got %d", len(r.Results))
		}
	}
}

func TestApplyThreshold_EmptyTargets(t *testing.T) {
	results := ApplyThreshold(map[string][]Result{}, DefaultThresholdOptions())
	if len(results) != 0 {
		t.Errorf("expected empty, got %d", len(results))
	}
}

func TestHasThresholdResults_True(t *testing.T) {
	r := []ThresholdResult{{Target: "x"}}
	if !HasThresholdResults(r) {
		t.Error("expected true")
	}
}

func TestHasThresholdResults_False(t *testing.T) {
	if HasThresholdResults(nil) {
		t.Error("expected false")
	}
}
