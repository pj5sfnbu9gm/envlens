package differ

import (
	"testing"
)

func TestIntersect_Empty(t *testing.T) {
	result := Intersect(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

func TestIntersect_SingleTarget(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"A": "1", "B": "2"},
	}
	result := Intersect(targets)
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
}

func TestIntersect_CommonKeys(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"A": "1", "B": "2", "C": "3"},
		"staging": {"A": "1", "B": "99"},
	}
	result := Intersect(targets)
	if len(result) != 2 {
		t.Fatalf("expected 2 intersecting keys, got %d", len(result))
	}
	if result[0].Key != "A" || result[1].Key != "B" {
		t.Fatalf("unexpected keys: %v", result)
	}
}

func TestIntersect_NoCommonKeys(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"A": "1"},
		"staging": {"B": "2"},
	}
	result := Intersect(targets)
	if len(result) != 0 {
		t.Fatalf("expected 0 common keys, got %d", len(result))
	}
}

func TestAllAgree_AllSame(t *testing.T) {
	results := []IntersectResult{
		{Key: "A", Values: map[string]string{"prod": "1", "staging": "1"}},
		{Key: "B", Values: map[string]string{"prod": "x", "staging": "y"}},
	}
	agreed := AllAgree(results)
	if len(agreed) != 1 || agreed[0] != "A" {
		t.Fatalf("expected [A], got %v", agreed)
	}
}

func TestAllAgree_NoneAgree(t *testing.T) {
	results := []IntersectResult{
		{Key: "A", Values: map[string]string{"prod": "1", "staging": "2"}},
	}
	if len(AllAgree(results)) != 0 {
		t.Fatal("expected no agreement")
	}
}
