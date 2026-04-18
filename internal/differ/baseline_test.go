package differ

import (
	"testing"
)

func baselineTargets() map[string]map[string]string {
	return map[string]map[string]string{
		"baseline": {"HOST": "localhost", "PORT": "5432", "DEBUG": "false"},
		"staging":  {"HOST": "staging.example.com", "PORT": "5432", "DEBUG": "true"},
		"prod":     {"HOST": "prod.example.com", "PORT": "5432"},
	}
}

func TestCompareToBaseline_ExcludesBaselineItself(t *testing.T) {
	results := CompareToBaseline(baselineTargets(), DefaultBaselineOptions())
	if _, ok := results["baseline"]; ok {
		t.Error("baseline target should not appear in results")
	}
}

func TestCompareToBaseline_DetectsChanges(t *testing.T) {
	results := CompareToBaseline(baselineTargets(), DefaultBaselineOptions())
	staging := results["staging"]
	if len(staging) == 0 {
		t.Fatal("expected diffs for staging")
	}
	found := false
	for _, r := range staging {
		if r.Key == "HOST" && r.Status == StatusChanged {
			found = true
		}
	}
	if !found {
		t.Error("expected HOST to be changed in staging")
	}
}

func TestCompareToBaseline_RemovedKey(t *testing.T) {
	results := CompareToBaseline(baselineTargets(), DefaultBaselineOptions())
	prod := results["prod"]
	found := false
	for _, r := range prod {
		if r.Key == "DEBUG" && r.Status == StatusRemoved {
			found = true
		}
	}
	if !found {
		t.Error("expected DEBUG to be removed in prod")
	}
}

func TestCompareToBaseline_IgnoreUnchangedFalse(t *testing.T) {
	opts := DefaultBaselineOptions()
	opts.IgnoreUnchanged = false
	results := CompareToBaseline(baselineTargets(), opts)
	staging := results["staging"]
	hasUnchanged := false
	for _, r := range staging {
		if r.Status == StatusUnchanged {
			hasUnchanged = true
		}
	}
	if !hasUnchanged {
		t.Error("expected unchanged keys when IgnoreUnchanged=false")
	}
}

func TestCompareToBaseline_MissingBaseline(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"HOST": "prod.example.com"},
	}
	results := CompareToBaseline(targets, DefaultBaselineOptions())
	prod := results["prod"]
	if len(prod) == 0 {
		t.Fatal("expected diffs when baseline is missing")
	}
	if prod[0].Status != StatusAdded {
		t.Errorf("expected Added status, got %s", prod[0].Status)
	}
}

func TestHasBaselineDifferences_True(t *testing.T) {
	results := CompareToBaseline(baselineTargets(), DefaultBaselineOptions())
	if !HasBaselineDifferences(results) {
		t.Error("expected differences")
	}
}

func TestHasBaselineDifferences_False(t *testing.T) {
	targets := map[string]map[string]string{
		"baseline": {"A": "1"},
		"prod":     {"A": "1"},
	}
	results := CompareToBaseline(targets, DefaultBaselineOptions())
	if HasBaselineDifferences(results) {
		t.Error("expected no differences")
	}
}
