package redactor

import (
	"testing"
)

func TestRedact_NoMatches(t *testing.T) {
	env := map[string]string{"APP_NAME": "envlens", "PORT": "8080"}
	r := Redact(env, DefaultOptions())
	if r.Redacted["APP_NAME"] != "envlens" {
		t.Errorf("expected value unchanged, got %s", r.Redacted["APP_NAME"])
	}
	if len(r.RedactedKeys) != 0 {
		t.Errorf("expected no redacted keys, got %v", r.RedactedKeys)
	}
}

func TestRedact_ByPrefix(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abc123", "APP_NAME": "envlens"}
	r := Redact(env, DefaultOptions())
	if r.Redacted["SECRET_KEY"] != "[REDACTED]" {
		t.Errorf("expected redacted, got %s", r.Redacted["SECRET_KEY"])
	}
	if r.Redacted["APP_NAME"] != "envlens" {
		t.Errorf("expected unchanged, got %s", r.Redacted["APP_NAME"])
	}
}

func TestRedact_ByExplicitKey(t *testing.T) {
	opts := DefaultOptions()
	opts.Keys = []string{"DB_PASS"}
	env := map[string]string{"DB_PASS": "hunter2", "DB_HOST": "localhost"}
	r := Redact(env, opts)
	if r.Redacted["DB_PASS"] != "[REDACTED]" {
		t.Errorf("expected redacted, got %s", r.Redacted["DB_PASS"])
	}
	if r.Redacted["DB_HOST"] != "localhost" {
		t.Errorf("expected unchanged, got %s", r.Redacted["DB_HOST"])
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	opts := DefaultOptions()
	opts.Placeholder = "***"
	env := map[string]string{"TOKEN_X": "secret"}
	r := Redact(env, opts)
	if r.Redacted["TOKEN_X"] != "***" {
		t.Errorf("expected ***, got %s", r.Redacted["TOKEN_X"])
	}
}

func TestRedact_OriginalUnmodified(t *testing.T) {
	env := map[string]string{"SECRET_FOO": "bar"}
	r := Redact(env, DefaultOptions())
	if r.Original["SECRET_FOO"] != "bar" {
		t.Errorf("original should be unmodified, got %s", r.Original["SECRET_FOO"])
	}
}

func TestRedact_EmptyEnv(t *testing.T) {
	r := Redact(map[string]string{}, DefaultOptions())
	if len(r.Redacted) != 0 {
		t.Errorf("expected empty redacted map")
	}
}
