package differ

import (
	"testing"
)

func sampleProjectionResults() []DiffResult {
	return []DiffResult{
		{Key: "APP_HOST", Status: StatusUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "APP_PORT", Status: StatusChanged, OldValue: "8080", NewValue: "9090"},
		{Key: "DB_HOST", Status: StatusAdded, OldValue: "", NewValue: "db.prod"},
		{Key: "DB_PASS", Status: StatusRemoved, OldValue: "secret", NewValue: ""},
		{Key: "LOG_LEVEL", Status: StatusUnchanged, OldValue: "info", NewValue: "info"},
	}
}

func TestProject_NoOptions(t *testing.T) {
	results := sampleProjectionResults()
	out := Project(results, DefaultProjectionOptions())
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestProject_ByExactKey(t *testing.T) {
	out := Project(sampleProjectionResults(), ProjectionOptions{Keys: []string{"APP_PORT", "LOG_LEVEL"}})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	if out[0].Key != "APP_PORT" || out[1].Key != "LOG_LEVEL" {
		t.Errorf("unexpected keys: %v", out)
	}
}

func TestProject_ByPrefix(t *testing.T) {
	out := Project(sampleProjectionResults(), ProjectionOptions{Prefixes: []string{"DB_"}})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestProject_Inverted(t *testing.T) {
	out := Project(sampleProjectionResults(), ProjectionOptions{Prefixes: []string{"DB_"}, Invert: true})
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
	for _, r := range out {
		if r.Key == "DB_HOST" || r.Key == "DB_PASS" {
			t.Errorf("DB_ key should have been excluded: %s", r.Key)
		}
	}
}

func TestHasProjectedChanges_True(t *testing.T) {
	out := Project(sampleProjectionResults(), ProjectionOptions{Prefixes: []string{"APP_"}})
	if !HasProjectedChanges(out) {
		t.Error("expected changes")
	}
}

func TestHasProjectedChanges_False(t *testing.T) {
	out := Project(sampleProjectionResults(), ProjectionOptions{Keys: []string{"APP_HOST", "LOG_LEVEL"}})
	if HasProjectedChanges(out) {
		t.Error("expected no changes")
	}
}
