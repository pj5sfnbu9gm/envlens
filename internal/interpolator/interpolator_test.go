package interpolator

import (
	"os"
	"testing"
)

func TestInterpolate_NoRefs(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{"BASE": "/app", "PATH": "${BASE}/bin"}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PATH"] != "/app/bin" {
		t.Errorf("expected /app/bin, got %q", out["PATH"])
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "URL": "http://$HOST:8080"}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %q", out["URL"])
	}
}

func TestInterpolate_FallbackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := map[string]string{"RESULT": "${OS_VAR}/suffix"}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["RESULT"] != "from-os/suffix" {
		t.Errorf("expected from-os/suffix, got %q", out["RESULT"])
	}
}

func TestInterpolate_MissingVar_NoFail(t *testing.T) {
	opts := Options{FallbackToOS: false, FailOnMissing: false}
	env := map[string]string{"VAL": "${MISSING}_end"}
	out, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAL"] != "_end" {
		t.Errorf("expected _end, got %q", out["VAL"])
	}
}

func TestInterpolate_MissingVar_FailOnMissing(t *testing.T) {
	opts := Options{FallbackToOS: false, FailOnMissing: true}
	env := map[string]string{"VAL": "${MISSING}"}
	_, err := Interpolate(env, opts)
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "${A} world"}
	original := map[string]string{"A": "hello", "B": "${A} world"}
	_, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range original {
		if env[k] != v {
			t.Errorf("input mutated: key %q changed to %q", k, env[k])
		}
	}
}
