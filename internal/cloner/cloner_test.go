package cloner_test

import (
	"testing"

	"github.com/user/envlens/internal/cloner"
)

var base = map[string]string{
	"APP_HOST": "localhost",
	"APP_PORT": "8080",
	"db_pass":  "secret",
	"DB_NAME":  "mydb",
}

func TestClone_NoTransform(t *testing.T) {
	out := cloner.Clone(base, cloner.DefaultOptions())
	if len(out) != len(base) {
		t.Fatalf("expected %d keys, got %d", len(base), len(out))
	}
	for k, v := range base {
		if out[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestClone_IsolatedCopy(t *testing.T) {
	out := cloner.Clone(base, cloner.DefaultOptions())
	out["APP_HOST"] = "changed"
	if base["APP_HOST"] == "changed" {
		t.Error("clone modified original map")
	}
}

func TestClone_KeyPrefix(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.KeyPrefix = "CLONE_"
	out := cloner.Clone(base, opts)
	if _, ok := out["CLONE_APP_HOST"]; !ok {
		t.Error("expected CLONE_APP_HOST key")
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("original key should not be present")
	}
}

func TestClone_UppercaseKeys(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.UppercaseKeys = true
	out := cloner.Clone(base, opts)
	if _, ok := out["DB_PASS"]; !ok {
		t.Error("expected DB_PASS after uppercase transform")
	}
	if _, ok := out["db_pass"]; ok {
		t.Error("lowercase key should not remain")
	}
}

func TestClone_FilterKeys(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.FilterKeys = []string{"APP_"}
	out := cloner.Clone(base, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
}

func TestClone_KeySuffix(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.KeySuffix = "_COPY"
	out := cloner.Clone(base, opts)
	if _, ok := out["APP_HOST_COPY"]; !ok {
		t.Error("expected APP_HOST_COPY")
	}
}

func TestClone_EmptyMap(t *testing.T) {
	out := cloner.Clone(map[string]string{}, cloner.DefaultOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d keys", len(out))
	}
}

func TestClone_NilMap(t *testing.T) {
	out := cloner.Clone(nil, cloner.DefaultOptions())
	if out == nil {
		t.Fatal("expected non-nil map returned for nil input")
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d keys", len(out))
	}
}
