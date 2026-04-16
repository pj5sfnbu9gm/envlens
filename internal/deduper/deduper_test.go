package deduper_test

import (
	"testing"

	"github.com/user/envlens/internal/deduper"
)

func TestDedupe_NoCaseFold_NoChange(t *testing.T) {
	env := map[string]string{"FOO": "1", "BAR": "2"}
	r := deduper.Dedupe(env, deduper.DefaultOptions())
	if len(r.Removed) != 0 {
		t.Errorf("expected no removals, got %v", r.Removed)
	}
	if r.Env["FOO"] != "1" || r.Env["BAR"] != "2" {
		t.Error("env values should be preserved")
	}
}

func TestDedupe_CaseFold_RemovesDuplicate(t *testing.T) {
	env := map[string]string{"foo": "lower", "FOO": "upper"}
	opts := deduper.Options{CaseFold: true}
	r := deduper.Dedupe(env, opts)
	if len(r.Removed) != 1 {
		t.Fatalf("expected 1 removal, got %d", len(r.Removed))
	}
	if len(r.Env) != 1 {
		t.Errorf("expected 1 key in output, got %d", len(r.Env))
	}
}

func TestDedupe_CaseFold_KeptMapping(t *testing.T) {
	env := map[string]string{"foo": "a", "FOO": "b"}
	r := deduper.Dedupe(env, deduper.Options{CaseFold: true})
	for dup, orig := range r.Kept {
		if _, ok := r.Env[orig]; !ok {
			t.Errorf("kept key %q (for dup %q) not present in output env", orig, dup)
		}
	}
}

func TestDedupe_CaseFold_MultipleGroups(t *testing.T) {
	env := map[string]string{
		"db_host": "a", "DB_HOST": "b", "Db_Host": "c",
		"PORT": "8080",
	}
	r := deduper.Dedupe(env, deduper.Options{CaseFold: true})
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removals, got %d: %v", len(r.Removed), r.Removed)
	}
	if len(r.Env) != 2 {
		t.Errorf("expected 2 keys in output, got %d", len(r.Env))
	}
}

func TestDedupe_EmptyEnv(t *testing.T) {
	r := deduper.Dedupe(map[string]string{}, deduper.DefaultOptions())
	if len(r.Env) != 0 || len(r.Removed) != 0 {
		t.Error("empty env should produce empty result")
	}
}
