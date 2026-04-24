package differ

import (
	"testing"
)

func buildLabelResults() []Result {
	return []Result{
		{Key: "APP_ENV", Status: StatusChanged, OldValue: "staging", NewValue: "production"},
		{Key: "DB_HOST", Status: StatusUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "NEW_KEY", Status: StatusAdded, OldValue: "", NewValue: "value1"},
		{Key: "OLD_KEY", Status: StatusRemoved, OldValue: "gone", NewValue: ""},
	}
}

func TestLabel_DefaultExcludesUnchanged(t *testing.T) {
	results := buildLabelResults()
	opts := DefaultLabelOptions()
	entries := Label(results, opts)
	for _, e := range entries {
		if e.Status == string(StatusUnchanged) {
			t.Errorf("expected unchanged entries to be excluded, got key %s", e.Key)
		}
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestLabel_IncludeUnchanged(t *testing.T) {
	results := buildLabelResults()
	opts := DefaultLabelOptions()
	opts.IncludeUnchanged = true
	entries := Label(results, opts)
	if len(entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(entries))
	}
}

func TestLabel_CustomLabel(t *testing.T) {
	results := buildLabelResults()
	opts := DefaultLabelOptions()
	opts.Labels = map[string]string{
		"APP_ENV": "environment-config",
		"DB_HOST": "database",
	}
	entries := Label(results, opts)
	for _, e := range entries {
		if e.Key == "APP_ENV" && e.Label != "environment-config" {
			t.Errorf("expected label 'environment-config', got %q", e.Label)
		}
	}
}

func TestLabel_DefaultLabel(t *testing.T) {
	results := buildLabelResults()
	opts := DefaultLabelOptions()
	opts.DefaultLabel = "generic"
	entries := Label(results, opts)
	for _, e := range entries {
		if e.Label != "generic" {
			t.Errorf("expected default label 'generic', got %q", e.Label)
		}
	}
}

func TestHasLabeledChanges_True(t *testing.T) {
	entries := []LabelEntry{
		{Key: "X", Status: string(StatusAdded), Label: "l"},
	}
	if !HasLabeledChanges(entries) {
		t.Error("expected HasLabeledChanges to return true")
	}
}

func TestHasLabeledChanges_False(t *testing.T) {
	entries := []LabelEntry{
		{Key: "X", Status: string(StatusUnchanged), Label: "l"},
	}
	if HasLabeledChanges(entries) {
		t.Error("expected HasLabeledChanges to return false")
	}
}
