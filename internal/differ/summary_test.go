package differ

import "testing"

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]Result{})
	if s.Total != 0 {
		t.Errorf("expected 0 total, got %d", s.Total)
	}
	if s.HasDifferences() {
		t.Error("expected no differences")
	}
}

func TestSummarize_Counts(t *testing.T) {
	results := []Result{
		{Status: StatusAdded},
		{Status: StatusAdded},
		{Status: StatusRemoved},
		{Status: StatusChanged},
		{Status: StatusUnchanged},
	}
	s := Summarize(results)
	if s.Total != 5 {
		t.Errorf("expected total 5, got %d", s.Total)
	}
	if s.Added != 2 {
		t.Errorf("expected 2 added, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", s.Changed)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", s.Unchanged)
	}
}

func TestSummarize_HasDifferences_True(t *testing.T) {
	results := []Result{{Status: StatusAdded}}
	if !Summarize(results).HasDifferences() {
		t.Error("expected HasDifferences to be true")
	}
}

func TestSummarize_HasDifferences_False(t *testing.T) {
	results := []Result{
		{Status: StatusUnchanged},
		{Status: StatusUnchanged},
	}
	if Summarize(results).HasDifferences() {
		t.Error("expected HasDifferences to be false")
	}
}
