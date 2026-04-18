package differ

import (
	"testing"
)

func TestMergeDiff_Empty(t *testing.T) {
	results := MergeDiff(nil, DefaultMergeOptions())
	if len(results) != 0 {
		t.Fatalf("expected empty, got %d", len(results))
	}
}

func TestMergeDiff_NoConflict(t *testing.T) {
	targets := map[string]map[string]string{
		"a": {"FOO": "bar"},
		"b": {"FOO": "bar"},
	}
	results := MergeDiff(targets, DefaultMergeOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Conflict {
		t.Error("expected no conflict")
	}
	if results[0].Value != "bar" {
		t.Errorf("expected bar, got %s", results[0].Value)
	}
}

func TestMergeDiff_ConflictDetected(t *testing.T) {
	targets := map[string]map[string]string{
		"a": {"FOO": "one"},
		"b": {"FOO": "two"},
	}
	results := MergeDiff(targets, DefaultMergeOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Conflict {
		t.Error("expected conflict")
	}
}

func TestMergeDiff_PreferFirst(t *testing.T) {
	opts := DefaultMergeOptions()
	opts.PreferFirst = true
	targets := map[string]map[string]string{
		"a": {"KEY": "first"},
		"b": {"KEY": "second"},
	}
	results := MergeDiff(targets, opts)
	if results[0].Value != "first" {
		t.Errorf("expected first, got %s", results[0].Value)
	}
}

func TestMergeDiff_LastWins(t *testing.T) {
	opts := DefaultMergeOptions()
	opts.PreferFirst = false
	targets := map[string]map[string]string{
		"a": {"KEY": "first"},
		"b": {"KEY": "second"},
	}
	results := MergeDiff(targets, opts)
	if results[0].Value != "second" {
		t.Errorf("expected second, got %s", results[0].Value)
	}
}

func TestMergeDiff_SkipConflicts(t *testing.T) {
	opts := DefaultMergeOptions()
	opts.SkipConflicts = true
	targets := map[string]map[string]string{
		"a": {"CONFLICT": "x", "SAFE": "ok"},
		"b": {"CONFLICT": "y", "SAFE": "ok"},
	}
	results := MergeDiff(targets, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result after skip, got %d", len(results))
	}
	if results[0].Key != "SAFE" {
		t.Errorf("expected SAFE key, got %s", results[0].Key)
	}
}

func TestHasMergeConflicts(t *testing.T) {
	results := []MergeResult{
		{Key: "A", Conflict: false},
		{Key: "B", Conflict: true},
	}
	if !HasMergeConflicts(results) {
		t.Error("expected conflicts")
	}
	if HasMergeConflicts(results[:1]) {
		t.Error("expected no conflicts")
	}
}
