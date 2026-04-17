package stringer_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/stringer"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "envlens",
		"DEBUG":    "true",
		"PORT":     "8080",
	}
}

func TestStringify_DefaultOptions(t *testing.T) {
	env := sampleEnv()
	opts := stringer.DefaultOptions()
	out := stringer.Stringify(env, opts)
	lines := strings.Split(out, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "APP_NAME=envlens" {
		t.Errorf("unexpected first line: %s", lines[0])
	}
}

func TestStringify_CustomSeparator(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	opts := stringer.DefaultOptions()
	opts.Separator = ": "
	out := stringer.Stringify(env, opts)
	if out != "KEY: val" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestStringify_QuoteValues(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	opts := stringer.DefaultOptions()
	opts.QuoteValues = true
	out := stringer.Stringify(env, opts)
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value in: %s", out)
	}
}

func TestStringify_EmptyEnv(t *testing.T) {
	out := stringer.Stringify(map[string]string{}, stringer.DefaultOptions())
	if out != "" {
		t.Errorf("expected empty string, got: %s", out)
	}
}

func TestToLines_Count(t *testing.T) {
	env := sampleEnv()
	lines := stringer.ToLines(env, stringer.DefaultOptions())
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestToLines_EmptyEnv(t *testing.T) {
	lines := stringer.ToLines(map[string]string{}, stringer.DefaultOptions())
	if len(lines) != 0 {
		t.Errorf("expected empty slice")
	}
}
