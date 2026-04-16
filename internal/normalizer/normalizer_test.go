package normalizer

import (
	"testing"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{"app_name": "envlens"}
	opts := DefaultOptions()
	out, results := Normalize(env, opts)

	if _, ok := out["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in output")
	}
	if len(results) != 1 || !results[0].Changed {
		t.Error("expected result to be marked as changed")
	}
}

func TestNormalize_TrimSpace(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	opts := DefaultOptions()
	out, _ := Normalize(env, opts)

	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestNormalize_ReplaceHyphens(t *testing.T) {
	env := map[string]string{"my-key": "val"}
	opts := DefaultOptions()
	out, _ := Normalize(env, opts)

	if _, ok := out["MY_KEY"]; !ok {
		t.Error("expected MY_KEY after hyphen replacement and uppercasing")
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	env := map[string]string{"PRESENT": "yes", "EMPTY": ""}
	opts := DefaultOptions()
	opts.RemoveEmpty = true
	out, results := Normalize(env, opts)

	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestNormalize_NoChange(t *testing.T) {
	env := map[string]string{"ALREADY_GOOD": "value"}
	opts := DefaultOptions()
	_, results := Normalize(env, opts)

	if len(results) != 1 {
		t.Fatal("expected 1 result")
	}
	if results[0].Changed {
		t.Error("expected result not to be marked as changed")
	}
}

func TestNormalize_EmptyEnv(t *testing.T) {
	out, results := Normalize(map[string]string{}, DefaultOptions())
	if len(out) != 0 || len(results) != 0 {
		t.Error("expected empty output for empty input")
	}
}
