package differ

import (
	"testing"
)

func TestFindUnique_Empty(t *testing.T) {
	result := FindUnique(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty, got %v", result)
	}
}

func TestFindUnique_SingleTarget(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"A": "1", "B": "2"},
	}
	results := FindUnique(targets)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	// All keys are unique when there's only one target.
	if len(results[0].Keys) != 2 {
		t.Errorf("expected 2 unique keys, got %v", results[0].Keys)
	}
}

func TestFindUnique_CommonKeysExcluded(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"SHARED": "1", "PROD_ONLY": "x"},
		"staging": {"SHARED": "1", "STAGE_ONLY": "y"},
	}
	results := FindUnique(targets)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	byTarget := map[string][]string{}
	for _, r := range results {
		byTarget[r.Target] = r.Keys
	}
	if len(byTarget["prod"]) != 1 || byTarget["prod"][0] != "PROD_ONLY" {
		t.Errorf("unexpected prod unique keys: %v", byTarget["prod"])
	}
	if len(byTarget["staging"]) != 1 || byTarget["staging"][0] != "STAGE_ONLY" {
		t.Errorf("unexpected staging unique keys: %v", byTarget["staging"])
	}
}

func TestFindUnique_NoUniqueKeys(t *testing.T) {
	targets := map[string]map[string]string{
		"a": {"X": "1"},
		"b": {"X": "2"},
	}
	results := FindUnique(targets)
	if HasUniqueKeys(results) {
		t.Error("expected no unique keys")
	}
}

func TestHasUniqueKeys_True(t *testing.T) {
	results := []UniqueResult{
		{Target: "a", Keys: []string{"ONLY_A"}},
		{Target: "b", Keys: nil},
	}
	if !HasUniqueKeys(results) {
		t.Error("expected HasUniqueKeys to be true")
	}
}

func TestFindUnique_SortedOutput(t *testing.T) {
	targets := map[string]map[string]string{
		"z": {"Z_KEY": "1"},
		"a": {"A_KEY": "2"},
	}
	results := FindUnique(targets)
	if results[0].Target != "a" || results[1].Target != "z" {
		t.Errorf("expected sorted target order, got %v, %v", results[0].Target, results[1].Target)
	}
}

func TestFindUnique_KeysSorted(t *testing.T) {
	targets := map[string]map[string]string{
		"only": {"ZEBRA": "1", "APPLE": "2", "MANGO": "3"},
	}
	results := FindUnique(targets)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	keys := results[0].Keys
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("expected sorted keys, got %v", keys)
			break
		}
	}
}
