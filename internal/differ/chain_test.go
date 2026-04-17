package differ

import (
	"testing"
)

func TestChain_Empty(t *testing.T) {
	result := Chain(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestChain_SingleTarget(t *testing.T) {
	targets := []NamedEnv{{Name: "a", Env: map[string]string{"K": "v"}}}
	result := Chain(targets)
	if len(result) != 0 {
		t.Errorf("expected 0 chain results, got %d", len(result))
	}
}

func TestChain_TwoTargets(t *testing.T) {
	targets := []NamedEnv{
		{Name: "dev", Env: map[string]string{"A": "1", "B": "old"}},
		{Name: "prod", Env: map[string]string{"A": "1", "B": "new"}},
	}
	chain := Chain(targets)
	if len(chain) != 1 {
		t.Fatalf("expected 1 step, got %d", len(chain))
	}
	if chain[0].From != "dev" || chain[0].To != "prod" {
		t.Errorf("unexpected from/to: %s/%s", chain[0].From, chain[0].To)
	}
	var changed int
	for _, r := range chain[0].Results {
		if r.Status == StatusChanged {
			changed++
		}
	}
	if changed != 1 {
		t.Errorf("expected 1 changed key, got %d", changed)
	}
}

func TestChain_ThreeTargets(t *testing.T) {
	targets := []NamedEnv{
		{Name: "dev", Env: map[string]string{"X": "1"}},
		{Name: "staging", Env: map[string]string{"X": "2"}},
		{Name: "prod", Env: map[string]string{"X": "3"}},
	}
	chain := Chain(targets)
	if len(chain) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(chain))
	}
	if chain[1].From != "staging" || chain[1].To != "prod" {
		t.Errorf("unexpected step 2: %s->%s", chain[1].From, chain[1].To)
	}
}

func TestHasChainChanges_True(t *testing.T) {
	chain := []ChainResult{{
		From: "a", To: "b",
		Results: []Result{{Key: "X", Status: StatusChanged}},
	}}
	if !HasChainChanges(chain) {
		t.Error("expected true")
	}
}

func TestHasChainChanges_False(t *testing.T) {
	chain := []ChainResult{{
		From: "a", To: "b",
		Results: []Result{{Key: "X", Status: StatusUnchanged}},
	}}
	if HasChainChanges(chain) {
		t.Error("expected false")
	}
}
