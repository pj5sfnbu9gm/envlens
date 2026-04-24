package differ

import (
	"testing"
)

func TestDetectDrift_Empty(t *testing.T) {
	r := DetectDrift("prod", nil)
	if len(r.Entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(r.Entries))
	}
}

func TestDetectDrift_SingleSnapshot(t *testing.T) {
	snap := map[string]string{"FOO": "bar"}
	r := DetectDrift("prod", []map[string]string{snap})
	if len(r.Entries) != 0 {
		t.Fatalf("expected no entries for single snapshot, got %d", len(r.Entries))
	}
}

func TestDetectDrift_Added(t *testing.T) {
	s1 := map[string]string{}
	s2 := map[string]string{"NEW_KEY": "value"}
	r := DetectDrift("staging", []map[string]string{s1, s2})
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Key != "NEW_KEY" || e.Status != "added" || e.To != "value" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDetectDrift_Removed(t *testing.T) {
	s1 := map[string]string{"OLD_KEY": "gone"}
	s2 := map[string]string{}
	r := DetectDrift("prod", []map[string]string{s1, s2})
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Key != "OLD_KEY" || e.Status != "removed" || e.From != "gone" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDetectDrift_Changed(t *testing.T) {
	s1 := map[string]string{"DB_HOST": "localhost"}
	s2 := map[string]string{"DB_HOST": "db.prod.internal"}
	r := DetectDrift("prod", []map[string]string{s1, s2})
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Status != "changed" || e.From != "localhost" || e.To != "db.prod.internal" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDetectDrift_Unchanged(t *testing.T) {
	s1 := map[string]string{"STABLE": "same"}
	s2 := map[string]string{"STABLE": "same"}
	r := DetectDrift("prod", []map[string]string{s1, s2})
	if len(r.Entries) != 0 {
		t.Fatalf("expected no entries for unchanged key, got %d", len(r.Entries))
	}
}

func TestDetectDrift_ThreeSnapshots(t *testing.T) {
	s1 := map[string]string{"X": "a"}
	s2 := map[string]string{"X": "b"}
	s3 := map[string]string{"X": "c"}
	r := DetectDrift("dev", []map[string]string{s1, s2, s3})
	// Two consecutive diffs: a->b and b->c
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestHasDrift_True(t *testing.T) {
	r := DriftReport{Entries: []DriftEntry{{Key: "K", Status: "added"}}}
	if !HasDrift(r) {
		t.Error("expected HasDrift to be true")
	}
}

func TestHasDrift_False(t *testing.T) {
	r := DriftReport{}
	if HasDrift(r) {
		t.Error("expected HasDrift to be false")
	}
}
