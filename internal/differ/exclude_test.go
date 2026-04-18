package differ

import (
	"testing"
)

func sampleExcludeResults() []Result {
	return []Result{
		{Key: "APP_HOST", Status: StatusChanged, OldValue: "a", NewValue: "b"},
		{Key: "DB_PASSWORD", Status: StatusChanged, OldValue: "x", NewValue: "y"},
		{Key: "DB_USER", Status: StatusAdded, NewValue: "admin"},
		{Key: "LOG_LEVEL", Status: StatusUnchanged, OldValue: "info", NewValue: "info"},
		{Key: "SECRET_KEY", Status: StatusRemoved, OldValue: "abc"},
	}
}

func TestExclude_NoOptions(t *testing.T) {
	results := sampleExcludeResults()
	out := Exclude(results, ExcludeOptions{})
	if len(out) != len(results) {
		t.Fatalf("expected %d, got %d", len(results), len(out))
	}
}

func TestExclude_ByExactKey(t *testing.T) {
	out := Exclude(sampleExcludeResults(), ExcludeOptions{Keys: []string{"LOG_LEVEL", "SECRET_KEY"}})
	for _, r := range out {
		if r.Key == "LOG_LEVEL" || r.Key == "SECRET_KEY" {
			t.Errorf("key %q should have been excluded", r.Key)
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestExclude_ByPrefix(t *testing.T) {
	out := Exclude(sampleExcludeResults(), ExcludeOptions{Prefixes: []string{"DB_"}})
	for _, r := range out {
		if r.Key == "DB_PASSWORD" || r.Key == "DB_USER" {
			t.Errorf("key %q should have been excluded", r.Key)
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestExclude_CombinedKeyAndPrefix(t *testing.T) {
	out := Exclude(sampleExcludeResults(), ExcludeOptions{
		Keys:     []string{"APP_HOST"},
		Prefixes: []string{"SECRET_"},
	})
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestHasExcluded_True(t *testing.T) {
	if !HasExcluded(sampleExcludeResults(), ExcludeOptions{Keys: []string{"LOG_LEVEL"}}) {
		t.Error("expected HasExcluded to return true")
	}
}

func TestHasExcluded_False(t *testing.T) {
	if HasExcluded(sampleExcludeResults(), ExcludeOptions{}) {
		t.Error("expected HasExcluded to return false")
	}
}
