package renamer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/renamer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"APP_SECRET":  "s3cr3t",
		"APP_VERSION": "1.0.0",
		"LOG_LEVEL":   "info",
	}
}

func TestRename_NoRules(t *testing.T) {
	env := baseEnv()
	out, err := renamer.Rename(env, renamer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(env) {
		t.Fatalf("expected %d keys, got %d", len(env), len(out))
	}
}

func TestRename_ExplicitRule(t *testing.T) {
	opts := renamer.Options{
		Rules: []renamer.Rule{{From: "DB_HOST", To: "DATABASE_HOST"}},
	}
	out, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("old key DB_HOST should be removed")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
}

func TestRename_MissingKeyIgnored(t *testing.T) {
	opts := renamer.Options{
		Rules: []renamer.Rule{{From: "NONEXISTENT", To: "NEW_KEY"}},
	}
	_, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("expected no error for missing key, got: %v", err)
	}
}

func TestRename_FailOnMissing(t *testing.T) {
	opts := renamer.Options{
		Rules:         []renamer.Rule{{From: "NONEXISTENT", To: "NEW_KEY"}},
		FailOnMissing: true,
	}
	_, err := renamer.Rename(baseEnv(), opts)
	if err == nil {
		t.Fatal("expected error for missing key with FailOnMissing")
	}
}

func TestRename_PrefixReplacement(t *testing.T) {
	opts := renamer.Options{
		OldPrefix: "APP_",
		NewPrefix: "SERVICE_",
	}
	out, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_SECRET"]; ok {
		t.Error("APP_SECRET should have been renamed")
	}
	if out["SERVICE_SECRET"] != "s3cr3t" {
		t.Errorf("expected SERVICE_SECRET=s3cr3t, got %q", out["SERVICE_SECRET"])
	}
	if out["SERVICE_VERSION"] != "1.0.0" {
		t.Errorf("expected SERVICE_VERSION=1.0.0, got %q", out["SERVICE_VERSION"])
	}
	// Keys without the prefix should be untouched.
	if out["DB_HOST"] != "localhost" {
		t.Error("DB_HOST should be unchanged")
	}
}

func TestRename_InvalidRule(t *testing.T) {
	opts := renamer.Options{
		Rules: []renamer.Rule{{From: "", To: "NEW_KEY"}},
	}
	_, err := renamer.Rename(baseEnv(), opts)
	if err == nil {
		t.Fatal("expected error for empty From field")
	}
}

func TestRename_OriginalUnmutated(t *testing.T) {
	env := baseEnv()
	opts := renamer.Options{
		Rules: []renamer.Rule{{From: "LOG_LEVEL", To: "LOGGING_LEVEL"}},
	}
	_, err := renamer.Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["LOG_LEVEL"]; !ok {
		t.Error("original map should not be mutated")
	}
}
