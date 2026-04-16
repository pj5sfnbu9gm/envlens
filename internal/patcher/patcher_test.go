package patcher_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/patcher"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_PASS":  "secret",
	}
}

func TestApply_SetNewKey(t *testing.T) {
	out, results, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "NEW_KEY", Value: "newval"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "newval" {
		t.Errorf("expected newval, got %q", out["NEW_KEY"])
	}
	if !results[0].Applied {
		t.Error("expected result to be applied")
	}
}

func TestApply_OverwriteExistingKey(t *testing.T) {
	out, _, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "APP_PORT", Value: "9090"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_PORT"] != "9090" {
		t.Errorf("expected 9090, got %q", out["APP_PORT"])
	}
}

func TestApply_UnsetKey(t *testing.T) {
	out, _, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpUnset, Key: "DB_PASS"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be removed")
	}
}

func TestApply_UnsetMissing_NoFail(t *testing.T) {
	_, results, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpUnset, Key: "NONEXISTENT"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Error("expected not applied for missing key")
	}
}

func TestApply_UnsetMissing_FailOnMissing(t *testing.T) {
	opts := patcher.Options{FailOnMissing: true}
	_, _, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpUnset, Key: "NONEXISTENT"},
	}, opts)
	if err == nil {
		t.Error("expected error for missing key with FailOnMissing")
	}
}

func TestApply_RenameKey(t *testing.T) {
	out, results, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Op: patcher.OpRename, Key: "APP_HOST", To: "SERVICE_HOST"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be removed")
	}
	if out["SERVICE_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", out["SERVICE_HOST"])
	}
	if !results[0].Applied {
		t.Error("expected result to be applied")
	}
}

func TestApply_UnknownOp(t *testing.T) {
	_, _, err := patcher.Apply(baseEnv(), []patcher.Patch{
		{Op: "invalid", Key: "APP_HOST"},
	}, patcher.DefaultOptions())
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	env := baseEnv()
	patcher.Apply(env, []patcher.Patch{ //nolint:errcheck
		{Op: patcher.OpSet, Key: "EXTRA", Value: "val"},
	}, patcher.DefaultOptions())
	if _, ok := env["EXTRA"]; ok {
		t.Error("original map should not be mutated")
	}
}
