package differ

import (
	"testing"
)

func buildSnapshots() []NamedEnv {
	return []NamedEnv{
		{Name: "t0", Env: map[string]string{"A": "1", "B": "hello", "C": "same"}},
		{Name: "t1", Env: map[string]string{"A": "2", "B": "hello", "D": "new"}},
		{Name: "t2", Env: map[string]string{"A": "3", "B": "world", "D": "new"}},
	}
}

func TestWindow_Empty(t *testing.T) {
	result := Window(nil, DefaultWindowOptions())
	if len(result) != 0 {
		t.Fatalf("expected no windows, got %d", len(result))
	}
}

func TestWindow_TooFewSnapshots(t *testing.T) {
	snaps := buildSnapshots()[:1]
	result := Window(snaps, DefaultWindowOptions())
	if len(result) != 0 {
		t.Fatalf("expected no windows for single snapshot, got %d", len(result))
	}
}

func TestWindow_TwoSnapshots_OneWindow(t *testing.T) {
	snaps := buildSnapshots()[:2]
	result := Window(snaps, DefaultWindowOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 window, got %d", len(result))
	}
	w := result[0]
	if w.Labels[0] != "t0" || w.Labels[1] != "t1" {
		t.Errorf("unexpected labels: %v", w.Labels)
	}
}

func TestWindow_ThreeSnapshots_TwoWindows(t *testing.T) {
	result := Window(buildSnapshots(), DefaultWindowOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(result))
	}
	if result[0].Labels[0] != "t0" || result[1].Labels[0] != "t1" {
		t.Errorf("unexpected window labels")
	}
}

func TestWindow_IgnoreUnchanged(t *testing.T) {
	snaps := buildSnapshots()[:2]
	opts := DefaultWindowOptions()
	opts.IgnoreUnchanged = true
	result := Window(snaps, opts)
	for _, r := range result[0].Results {
		if r.Status == StatusUnchanged {
			t.Errorf("expected unchanged entries to be filtered, got key %s", r.Key)
		}
	}
}

func TestWindow_IncludeUnchanged(t *testing.T) {
	snaps := buildSnapshots()[:2]
	opts := DefaultWindowOptions()
	opts.IgnoreUnchanged = false
	result := Window(snaps, opts)
	found := false
	for _, r := range result[0].Results {
		if r.Status == StatusUnchanged {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one unchanged entry when IgnoreUnchanged=false")
	}
}

func TestHasWindowChanges_True(t *testing.T) {
	result := Window(buildSnapshots(), DefaultWindowOptions())
	if !HasWindowChanges(result) {
		t.Error("expected HasWindowChanges to return true")
	}
}

func TestHasWindowChanges_False(t *testing.T) {
	snaps := []NamedEnv{
		{Name: "a", Env: map[string]string{"X": "1"}},
		{Name: "b", Env: map[string]string{"X": "1"}},
	}
	opts := DefaultWindowOptions()
	opts.IgnoreUnchanged = false
	result := Window(snaps, opts)
	if HasWindowChanges(result) {
		t.Error("expected HasWindowChanges to return false for identical snapshots")
	}
}
