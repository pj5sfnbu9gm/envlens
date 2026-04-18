package freezer

import (
	"testing"
)

func TestFreeze_Basic(t *testing.T) {
	src := map[string]string{"HOST": "localhost", "PORT": "8080"}
	f, err := Freeze(src, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Len() != 2 {
		t.Errorf("expected 2 keys, got %d", f.Len())
	}
}

func TestFreeze_NilSource(t *testing.T) {
	_, err := Freeze(nil, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestFreeze_EmptyKeyRejected(t *testing.T) {
	src := map[string]string{"": "value"}
	_, err := Freeze(src, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestFreeze_AllowEmptyFalse(t *testing.T) {
	src := map[string]string{"KEY": ""}
	opts := Options{AllowEmpty: false}
	_, err := Freeze(src, opts)
	if err == nil {
		t.Fatal("expected error for empty value when AllowEmpty=false")
	}
}

func TestFreeze_IsolatesSource(t *testing.T) {
	src := map[string]string{"A": "1"}
	f, _ := Freeze(src, DefaultOptions())
	src["A"] = "mutated"
	if v, _ := f.Get("A"); v != "1" {
		t.Errorf("freeze should isolate from source mutations")
	}
}

func TestToMap_IsolatesFrozen(t *testing.T) {
	f, _ := Freeze(map[string]string{"X": "10"}, DefaultOptions())
	m := f.ToMap()
	m["X"] = "changed"
	if v, _ := f.Get("X"); v != "10" {
		t.Errorf("ToMap should return independent copy")
	}
}

func TestDiff_NoDifferences(t *testing.T) {
	a, _ := Freeze(map[string]string{"K": "v"}, DefaultOptions())
	b, _ := Freeze(map[string]string{"K": "v"}, DefaultOptions())
	if len(Diff(a, b)) != 0 {
		t.Error("expected no differences")
	}
}

func TestDiff_DetectsChange(t *testing.T) {
	a, _ := Freeze(map[string]string{"K": "v1"}, DefaultOptions())
	b, _ := Freeze(map[string]string{"K": "v2"}, DefaultOptions())
	if len(Diff(a, b)) != 1 {
		t.Error("expected one difference")
	}
}

func TestDiff_DetectsMissingKey(t *testing.T) {
	a, _ := Freeze(map[string]string{"A": "1", "B": "2"}, DefaultOptions())
	b, _ := Freeze(map[string]string{"A": "1"}, DefaultOptions())
	if len(Diff(a, b)) != 1 {
		t.Errorf("expected 1 diff for missing key, got %d", len(Diff(a, b)))
	}
}
