package sorter_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/sorter"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"APP_NAME":    "envlens",
		"DB_PORT":     "5432",
		"APP_VERSION": "1.0.0",
		"LOG_LEVEL":   "info",
	}
}

func TestSort_Ascending(t *testing.T) {
	_, keys := sorter.Sort(sampleEnv(), sorter.DefaultOptions())
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Errorf("expected ascending order, got %q before %q", keys[i-1], keys[i])
		}
	}
}

func TestSort_Descending(t *testing.T) {
	opts := sorter.Options{Order: sorter.Descending}
	_, keys := sorter.Sort(sampleEnv(), opts)
	for i := 1; i < len(keys); i++ {
		if keys[i-1] < keys[i] {
			t.Errorf("expected descending order, got %q before %q", keys[i-1], keys[i])
		}
	}
}

func TestSort_GroupByPrefix(t *testing.T) {
	opts := sorter.Options{Order: sorter.Ascending, GroupByPrefix: true}
	_, keys := sorter.Sort(sampleEnv(), opts)

	// All APP_ keys should appear before DB_ keys which appear before LOG_
	prefixOrder := map[string]int{}
	cursor := 0
	for _, k := range keys {
		pfx := prefixOf(k)
		if _, seen := prefixOrder[pfx]; !seen {
			prefixOrder[pfx] = cursor
			cursor++
		}
	}
	if prefixOrder["APP"] >= prefixOrder["DB"] {
		t.Errorf("expected APP prefix before DB prefix")
	}
	if prefixOrder["DB"] >= prefixOrder["LOG"] {
		t.Errorf("expected DB prefix before LOG prefix")
	}
}

func TestSort_PreservesValues(t *testing.T) {
	env := sampleEnv()
	out, _ := sorter.Sort(env, sorter.DefaultOptions())
	for k, v := range env {
		if out[k] != v {
			t.Errorf("value mismatch for key %q: got %q, want %q", k, out[k], v)
		}
	}
}

func TestSort_EmptyMap(t *testing.T) {
	out, keys := sorter.Sort(map[string]string{}, sorter.DefaultOptions())
	if len(out) != 0 || len(keys) != 0 {
		t.Errorf("expected empty results for empty input")
	}
}

func prefixOf(key string) string {
	for i, ch := range key {
		if ch == '_' && i > 0 {
			return key[:i]
		}
	}
	return key
}
