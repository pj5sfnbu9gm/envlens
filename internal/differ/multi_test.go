package differ

import (
	"testing"
)

func TestMultiDiff_NoChanges(t *testing.T) {
	baseline := map[string]string{"A": "1", "B": "2"}
	targets := map[string]map[string]string{
		"prod": {"A": "1", "B": "2"},
	}
	diffs, err := MultiDiff(baseline, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 target diff, got %d", len(diffs))
	}
	if HasAnyChanges(diffs) {
		t.Error("expected no changes")
	}
}

func TestMultiDiff_DetectsChanges(t *testing.T) {
	baseline := map[string]string{"A": "1", "B": "2"}
	targets := map[string]map[string]string{
		"staging": {"A": "1", "B": "99"},
	}
	diffs, err := MultiDiff(baseline, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasAnyChanges(diffs) {
		t.Error("expected changes")
	}
}

func TestMultiDiff_MultipleTargets(t *testing.T) {
	baseline := map[string]string{"X": "a"}
	targets := map[string]map[string]string{
		"alpha": {"X": "a"},
		"beta":  {"X": "b"},
		"gamma": {"X": "a", "Y": "new"},
	}
	diffs, err := MultiDiff(baseline, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(diffs))
	}
	// diffs should be sorted by target name: alpha, beta, gamma
	if diffs[0].Target != "alpha" || diffs[1].Target != "beta" || diffs[2].Target != "gamma" {
		t.Errorf("unexpected order: %v", []string{diffs[0].Target, diffs[1].Target, diffs[2].Target})
	}
}

func TestMultiDiff_NilBaseline(t *testing.T) {
	_, err := MultiDiff(nil, map[string]map[string]string{})
	if err == nil {
		t.Error("expected error for nil baseline")
	}
}

func TestMultiDiff_EmptyTargets(t *testing.T) {
	baseline := map[string]string{"A": "1"}
	diffs, err := MultiDiff(baseline, map[string]map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(diffs))
	}
	if HasAnyChanges(diffs) {
		t.Error("expected no changes for empty targets")
	}
}
