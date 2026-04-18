package differ

import (
	"testing"
)

func TestFindOverlap_Empty(t *testing.T) {
	results := FindOverlap(nil)
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestFindOverlap_AllPresent(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"HOST": "a", "PORT": "80"},
		"staging": {"HOST": "b", "PORT": "8080"},
	}
	results := FindOverlap(targets)
	if len(results) != 0 {
		t.Fatalf("expected no overlap results, got %d", len(results))
	}
}

func TestFindOverlap_MissingInOne(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"HOST": "a", "PORT": "80", "SECRET": "x"},
		"staging": {"HOST": "b", "PORT": "8080"},
	}
	results := FindOverlap(targets)
	if len(results) != 1 {
		t.Fatalf("expected 1 overlap result, got %d", len(results))
	}
	if results[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %s", results[0].Key)
	}
	if len(results[0].PresentIn) != 1 || results[0].PresentIn[0] != "prod" {
		t.Errorf("unexpected PresentIn: %v", results[0].PresentIn)
	}
	if len(results[0].AbsentFrom) != 1 || results[0].AbsentFrom[0] != "staging" {
		t.Errorf("unexpected AbsentFrom: %v", results[0].AbsentFrom)
	}
}

func TestFindOverlap_MultipleTargets(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"HOST": "a", "ONLY_PROD": "1"},
		"staging": {"HOST": "b", "ONLY_STAGING": "2"},
		"dev":     {"HOST": "c"},
	}
	results := FindOverlap(targets)
	if len(results) != 2 {
		t.Fatalf("expected 2 overlap results, got %d", len(results))
	}
	keys := map[string]bool{}
	for _, r := range results {
		keys[r.Key] = true
	}
	if !keys["ONLY_PROD"] || !keys["ONLY_STAGING"] {
		t.Errorf("unexpected keys in overlap: %v", keys)
	}
}

func TestHasOverlap_True(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"HOST": "a", "EXTRA": "1"},
		"staging": {"HOST": "b"},
	}
	if !HasOverlap(targets) {
		t.Error("expected HasOverlap to return true")
	}
}

func TestHasOverlap_False(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"HOST": "a"},
		"staging": {"HOST": "b"},
	}
	if HasOverlap(targets) {
		t.Error("expected HasOverlap to return false")
	}
}
