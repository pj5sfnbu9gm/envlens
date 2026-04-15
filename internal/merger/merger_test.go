package merger

import (
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	sources := map[string]map[string]string{
		"base": {"APP_ENV": "production", "PORT": "8080"},
		"extra": {"DEBUG": "false", "TIMEOUT": "30"},
	}
	result, err := Merge(sources, []string{"base", "extra"}, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(result.Conflicts))
	}
	if result.Env["APP_ENV"] != "production" || result.Env["DEBUG"] != "false" {
		t.Errorf("unexpected env values: %v", result.Env)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	sources := map[string]map[string]string{
		"base":     {"PORT": "8080"},
		"override": {"PORT": "9090"},
	}
	result, err := Merge(sources, []string{"base", "override"}, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", result.Env["PORT"])
	}
	if len(result.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(result.Conflicts))
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	sources := map[string]map[string]string{
		"base":     {"PORT": "8080"},
		"override": {"PORT": "9090"},
	}
	result, err := Merge(sources, []string{"base", "override"}, StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["PORT"] != "9090" {
		t.Errorf("expected PORT=9090, got %s", result.Env["PORT"])
	}
	if len(result.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(result.Conflicts))
	}
}

func TestMerge_StrategyError(t *testing.T) {
	sources := map[string]map[string]string{
		"base":     {"PORT": "8080"},
		"override": {"PORT": "9090"},
	}
	_, err := Merge(sources, []string{"base", "override"}, StrategyError)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestMerge_MissingSourceInOrder(t *testing.T) {
	sources := map[string]map[string]string{
		"base": {"APP_ENV": "staging"},
	}
	result, err := Merge(sources, []string{"base", "nonexistent"}, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %s", result.Env["APP_ENV"])
	}
}

func TestMerge_EmptySources(t *testing.T) {
	result, err := Merge(map[string]map[string]string{}, []string{}, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Env) != 0 {
		t.Errorf("expected empty env, got %v", result.Env)
	}
}
