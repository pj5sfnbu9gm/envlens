package filter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/filter"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.internal",
		"DB_PASSWORD": "secret",
		"LOG_LEVEL":   "info",
	}
}

func TestApply_NoOptions(t *testing.T) {
	env := sampleEnv()
	result := filter.Apply(env, filter.Options{})
	if len(result) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(result))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	result := filter.Apply(sampleEnv(), filter.Options{
		Prefixes: []string{"APP_"},
	})
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("did not expect DB_HOST in result")
	}
}

func TestApply_ContainsFilter(t *testing.T) {
	result := filter.Apply(sampleEnv(), filter.Options{
		Contains: "host",
	})
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	result := filter.Apply(sampleEnv(), filter.Options{
		ExcludeKeys: []string{"DB_PASSWORD", "LOG_LEVEL"},
	})
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("did not expect DB_PASSWORD in result")
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	result := filter.Apply(sampleEnv(), filter.Options{
		Prefixes:    []string{"DB_"},
		ExcludeKeys: []string{"DB_PASSWORD"},
	})
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	result := filter.Apply(map[string]string{}, filter.Options{
		Prefixes: []string{"APP_"},
	})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}
