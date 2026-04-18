package differ

import (
	"testing"
)

func TestPivotDiff_Empty(t *testing.T) {
	result := PivotDiff(nil, DefaultPivotOptions())
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestPivotDiff_UniformExcludedByDefault(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"HOST": "example.com", "PORT": "8080"},
		"dev":  {"HOST": "example.com", "PORT": "8080"},
	}
	result := PivotDiff(targets, DefaultPivotOptions())
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestPivotDiff_DetectsDifference(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"HOST": "prod.example.com", "PORT": "443"},
		"dev":  {"HOST": "localhost", "PORT": "443"},
	}
	result := PivotDiff(targets, DefaultPivotOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", result[0].Key)
	}
	if result[0].Uniform {
		t.Error("expected non-uniform")
	}
}

func TestPivotDiff_IncludeUnchanged(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"HOST": "prod.example.com", "PORT": "443"},
		"dev":  {"HOST": "localhost", "PORT": "443"},
	}
	opts := DefaultPivotOptions()
	opts.IncludeUnchanged = true
	result := PivotDiff(targets, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestPivotDiff_MissingKeyInTarget(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"SECRET": "abc"},
		"dev":  {},
	}
	result := PivotDiff(targets, DefaultPivotOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Values["dev"] != "" {
		t.Errorf("expected empty string for missing key")
	}
}

func TestHasPivotDifferences_True(t *testing.T) {
	entries := []PivotEntry{{Key: "X", Uniform: false}}
	if !HasPivotDifferences(entries) {
		t.Error("expected true")
	}
}

func TestHasPivotDifferences_False(t *testing.T) {
	entries := []PivotEntry{{Key: "X", Uniform: true}}
	if HasPivotDifferences(entries) {
		t.Error("expected false")
	}
}
