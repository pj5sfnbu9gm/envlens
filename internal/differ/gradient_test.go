package differ

import (
	"testing"
)

func buildNamedEnvs() []NamedEnv {
	return []NamedEnv{
		{Name: "dev", Env: map[string]string{"APP_PORT": "3000", "LOG_LEVEL": "debug", "SHARED": "x"}},
		{Name: "staging", Env: map[string]string{"APP_PORT": "4000", "LOG_LEVEL": "info", "SHARED": "x"}},
		{Name: "prod", Env: map[string]string{"APP_PORT": "8080", "LOG_LEVEL": "warn", "SHARED": "x"}},
	}
}

func TestGradient_Empty(t *testing.T) {
	result := Gradient(nil, DefaultGradientOptions())
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestGradient_CountsChanges(t *testing.T) {
	result := Gradient(buildNamedEnvs(), DefaultGradientOptions())
	if len(result) == 0 {
		t.Fatal("expected entries")
	}
	for _, e := range result {
		if e.Key == "APP_PORT" && e.Changes != 2 {
			t.Errorf("APP_PORT: expected 2 changes, got %d", e.Changes)
		}
	}
}

func TestGradient_StableKeyExcludedByDefault(t *testing.T) {
	result := Gradient(buildNamedEnvs(), DefaultGradientOptions())
	for _, e := range result {
		if e.Key == "SHARED" {
			t.Error("SHARED should be excluded (stable)")
		}
	}
}

func TestGradient_IncludeStable(t *testing.T) {
	opts := DefaultGradientOptions()
	opts.IncludeStable = true
	opts.MinChanges = 0
	result := Gradient(buildNamedEnvs(), opts)
	found := false
	for _, e := range result {
		if e.Key == "SHARED" {
			found = true
			if e.Direction != "stable" {
				t.Errorf("expected stable, got %s", e.Direction)
			}
		}
	}
	if !found {
		t.Error("expected SHARED to be present")
	}
}

func TestGradient_StepsMatchTargetOrder(t *testing.T) {
	result := Gradient(buildNamedEnvs(), DefaultGradientOptions())
	if len(result) == 0 {
		t.Fatal("no entries")
	}
	expected := []string{"dev", "staging", "prod"}
	for _, e := range result {
		for i, s := range e.Steps {
			if s != expected[i] {
				t.Errorf("step %d: expected %s, got %s", i, expected[i], s)
			}
		}
	}
}

func TestGradient_ValuesLength(t *testing.T) {
	result := Gradient(buildNamedEnvs(), DefaultGradientOptions())
	for _, e := range result {
		if len(e.Values) != 3 {
			t.Errorf("%s: expected 3 values, got %d", e.Key, len(e.Values))
		}
	}
}

func TestHasGradientChanges_True(t *testing.T) {
	entries := []GradientEntry{{Key: "X", Changes: 1}}
	if !HasGradientChanges(entries) {
		t.Error("expected true")
	}
}

func TestHasGradientChanges_False(t *testing.T) {
	entries := []GradientEntry{{Key: "X", Changes: 0}}
	if HasGradientChanges(entries) {
		t.Error("expected false")
	}
}
