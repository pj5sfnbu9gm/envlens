package tagger

import (
	"testing"
)

func TestTag_NoRules(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	results := Tag(env, DefaultOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if len(r.Tags) != 0 {
			t.Errorf("expected no tags for %s, got %v", r.Key, r.Tags)
		}
	}
}

func TestTag_PrefixRule(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_ENV": "prod"}
	opts := DefaultOptions()
	opts.PrefixTags["DB_"] = []string{"database"}
	results := Tag(env, opts)
	tagged := map[string][]string{}
	for _, r := range results {
		tagged[r.Key] = r.Tags
	}
	if len(tagged["DB_HOST"]) == 0 || tagged["DB_HOST"][0] != "database" {
		t.Errorf("expected DB_HOST to have tag 'database'")
	}
	if len(tagged["APP_ENV"]) != 0 {
		t.Errorf("expected APP_ENV to have no tags")
	}
}

func TestTag_ExplicitRule(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abc123"}
	opts := DefaultOptions()
	opts.ExplicitTags["SECRET_KEY"] = []string{"sensitive", "secret"}
	results := Tag(env, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if len(results[0].Tags) != 2 {
		t.Errorf("expected 2 tags, got %v", results[0].Tags)
	}
}

func TestTag_NoDuplicateTags(t *testing.T) {
	env := map[string]string{"DB_SECRET": "pass"}
	opts := DefaultOptions()
	opts.PrefixTags["DB_"] = []string{"sensitive"}
	opts.ExplicitTags["DB_SECRET"] = []string{"sensitive", "secret"}
	results := Tag(env, opts)
	counts := map[string]int{}
	for _, tag := range results[0].Tags {
		counts[tag]++
	}
	if counts["sensitive"] != 1 {
		t.Errorf("expected 'sensitive' exactly once, got %d", counts["sensitive"])
	}
}

func TestTag_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	results := Tag(env, DefaultOptions())
	keys := []string{results[0].Key, results[1].Key, results[2].Key}
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}
