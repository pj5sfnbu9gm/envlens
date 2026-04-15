package trimmer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/trimmer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "envlens",
		"APP_VERSION": "",
		"DEBUG":       "   ",
		"SECRET_KEY":  "abc123",
		"TMP_FILE":    "/tmp/foo",
		"LOG_LEVEL":   "info",
	}
}

func TestTrim_RemoveEmpty(t *testing.T) {
	env := baseEnv()
	opts := trimmer.Options{RemoveEmpty: true}
	out := trimmer.Trim(env, opts)
	if _, ok := out["APP_VERSION"]; ok {
		t.Error("expected APP_VERSION to be removed (empty value)")
	}
	if _, ok := out["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to be retained")
	}
}

func TestTrim_RemoveBlank(t *testing.T) {
	env := baseEnv()
	opts := trimmer.Options{RemoveBlank: true}
	out := trimmer.Trim(env, opts)
	if _, ok := out["DEBUG"]; ok {
		t.Error("expected DEBUG to be removed (blank value)")
	}
	if _, ok := out["APP_VERSION"]; !ok {
		t.Error("expected APP_VERSION to be retained when RemoveEmpty is false")
	}
}

func TestTrim_RemoveExactKeys(t *testing.T) {
	env := baseEnv()
	opts := trimmer.Options{RemoveKeys: []string{"SECRET_KEY", "LOG_LEVEL"}}
	out := trimmer.Trim(env, opts)
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("expected SECRET_KEY to be removed")
	}
	if _, ok := out["LOG_LEVEL"]; ok {
		t.Error("expected LOG_LEVEL to be removed")
	}
	if _, ok := out["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to be retained")
	}
}

func TestTrim_RemovePrefixes(t *testing.T) {
	env := baseEnv()
	opts := trimmer.Options{RemovePrefixes: []string{"TMP_", "APP_"}}
	out := trimmer.Trim(env, opts)
	if _, ok := out["TMP_FILE"]; ok {
		t.Error("expected TMP_FILE to be removed")
	}
	if _, ok := out["APP_NAME"]; ok {
		t.Error("expected APP_NAME to be removed")
	}
	if _, ok := out["SECRET_KEY"]; !ok {
		t.Error("expected SECRET_KEY to be retained")
	}
}

func TestTrim_DefaultOptions(t *testing.T) {
	env := baseEnv()
	opts := trimmer.DefaultOptions()
	out := trimmer.Trim(env, opts)
	if _, ok := out["APP_VERSION"]; ok {
		t.Error("expected APP_VERSION removed by default options")
	}
	if _, ok := out["DEBUG"]; ok {
		t.Error("expected DEBUG removed by default options")
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	env := baseEnv()
	opts := trimmer.DefaultOptions()
	_ = trimmer.Trim(env, opts)
	if _, ok := env["APP_VERSION"]; !ok {
		t.Error("original map should not be mutated")
	}
}
