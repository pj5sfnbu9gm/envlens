package differ

import (
	"testing"
)

func TestScope_NoOptions_ReturnsAll(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "1", "B": "3"}
	res := Scope(base, target, ScopeOptions{})
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestScope_ByKey(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	target := map[string]string{"A": "9", "B": "2", "C": "7"}
	res := Scope(base, target, ScopeOptions{Keys: []string{"A"}})
	if len(res) != 1 || res[0].Key != "A" {
		t.Fatalf("expected only key A, got %+v", res)
	}
}

func TestScope_ByPrefix(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_ENV": "dev"}
	target := map[string]string{"DB_HOST": "prod", "DB_PORT": "5432", "APP_ENV": "prod"}
	res := Scope(base, target, ScopeOptions{Prefixes: []string{"DB_"}})
	if len(res) != 2 {
		t.Fatalf("expected 2 DB_ results, got %d", len(res))
	}
	for _, r := range res {
		if len(r.Key) < 3 || r.Key[:3] != "DB_" {
			t.Errorf("unexpected key %s", r.Key)
		}
	}
}

func TestScope_IgnoreUnchanged(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "1", "B": "9"}
	res := Scope(base, target, ScopeOptions{IgnoreUnchanged: true})
	if len(res) != 1 || res[0].Key != "B" {
		t.Fatalf("expected only changed key B, got %+v", res)
	}
}

func TestHasScopeChanges_True(t *testing.T) {
	results := []ScopeResult{{Key: "X", Status: "changed", Old: "a", New: "b"}}
	if !HasScopeChanges(results) {
		t.Error("expected HasScopeChanges to be true")
	}
}

func TestHasScopeChanges_False(t *testing.T) {
	results := []ScopeResult{{Key: "X", Status: "unchanged", Old: "a", New: "a"}}
	if HasScopeChanges(results) {
		t.Error("expected HasScopeChanges to be false")
	}
}
