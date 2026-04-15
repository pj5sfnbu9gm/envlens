package comparator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlens/internal/comparator"
	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/resolver"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestCompareAll_NoChanges(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	other := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	targets := []resolver.Target{
		{Name: "prod", Path: base},
		{Name: "staging", Path: other},
	}
	diffs, err := comparator.CompareAll("prod", targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comparator.HasChanges(diffs) {
		t.Error("expected no changes")
	}
}

func TestCompareAll_DetectsChanges(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nONLY_BASE=1\n")
	other := writeTempEnv(t, "FOO=changed\nONLY_OTHER=2\n")

	targets := []resolver.Target{
		{Name: "prod", Path: base},
		{Name: "staging", Path: other},
	}
	diffs, err := comparator.CompareAll("prod", targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !comparator.HasChanges(diffs) {
		t.Error("expected changes to be detected")
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	statuses := make(map[string]differ.Status)
	for _, r := range diffs[0].Results {
		statuses[r.Key] = r.Status
	}
	if statuses["FOO"] != differ.Changed {
		t.Errorf("expected FOO to be Changed, got %v", statuses["FOO"])
	}
	if statuses["ONLY_BASE"] != differ.Removed {
		t.Errorf("expected ONLY_BASE to be Removed, got %v", statuses["ONLY_BASE"])
	}
	if statuses["ONLY_OTHER"] != differ.Added {
		t.Errorf("expected ONLY_OTHER to be Added, got %v", statuses["ONLY_OTHER"])
	}
}

func TestCompareAll_MissingBaseline(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	targets := []resolver.Target{{Name: "staging", Path: path}}
	_, err := comparator.CompareAll("prod", targets)
	if err == nil {
		t.Error("expected error for missing baseline")
	}
}

func TestCompareAll_SkipsBaseline(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\n")
	targets := []resolver.Target{{Name: "prod", Path: base}}
	diffs, err := comparator.CompareAll("prod", targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(diffs))
	}
}
