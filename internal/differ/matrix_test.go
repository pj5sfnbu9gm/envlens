package differ

import (
	"testing"
)

func TestMatrix_Empty(t *testing.T) {
	entries := Matrix(map[string]map[string]string{})
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestMatrix_TwoTargets_PairsGenerated(t *testing.T) {
	targets := map[string]map[string]string{
		"dev":  {"A": "1"},
		"prod": {"A": "2"},
	}
	entries := Matrix(targets)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries (dev→prod, prod→dev), got %d", len(entries))
	}
}

func TestMatrix_DetectsChanges(t *testing.T) {
	targets := map[string]map[string]string{
		"dev":  {"A": "1", "B": "hello"},
		"prod": {"A": "2", "B": "hello"},
	}
	entries := Matrix(targets)
	if !HasMatrixChanges(entries) {
		t.Fatal("expected changes to be detected")
	}
}

func TestMatrix_NoChanges(t *testing.T) {
	targets := map[string]map[string]string{
		"dev":  {"A": "1"},
		"prod": {"A": "1"},
	}
	entries := Matrix(targets)
	if HasMatrixChanges(entries) {
		t.Fatal("expected no changes")
	}
}

func TestMatrix_ThreeTargets_SixPairs(t *testing.T) {
	targets := map[string]map[string]string{
		"a": {"X": "1"},
		"b": {"X": "1"},
		"c": {"X": "1"},
	}
	entries := Matrix(targets)
	if len(entries) != 6 {
		t.Fatalf("expected 6 pairs, got %d", len(entries))
	}
}

func TestMatrix_AddedRemovedKeys(t *testing.T) {
	targets := map[string]map[string]string{
		"dev":  {"A": "1", "B": "2"},
		"prod": {"A": "1", "C": "3"},
	}
	entries := Matrix(targets)
	var devToProd MatrixEntry
	for _, e := range entries {
		if e.From == "dev" && e.To == "prod" {
			devToProd = e
		}
	}
	statuses := map[string]string{}
	for _, r := range devToProd.Results {
		statuses[r.Key] = string(r.Status)
	}
	if statuses["B"] != string(StatusRemoved) {
		t.Errorf("expected B removed, got %s", statuses["B"])
	}
	if statuses["C"] != string(StatusAdded) {
		t.Errorf("expected C added, got %s", statuses["C"])
	}
}
