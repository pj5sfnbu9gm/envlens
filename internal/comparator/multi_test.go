package comparator

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
	return p
}

func TestMultiCompareAll_NoChanges(t *testing.T) {
	dir := t.TempDir()
	base := writeEnvFile(t, dir, "base.env", "FOO=bar\nBAZ=qux\n")
	tgt := writeEnvFile(t, dir, "tgt.env", "FOO=bar\nBAZ=qux\n")

	results, err := MultiCompareAll(MultiCompareOptions{
		Baseline: base,
		Targets:  map[string]string{"prod": tgt},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if AnyTargetHasChanges(results) {
		t.Error("expected no changes")
	}
}

func TestMultiCompareAll_DetectsChanges(t *testing.T) {
	dir := t.TempDir()
	base := writeEnvFile(t, dir, "base.env", "FOO=bar\n")
	tgt := writeEnvFile(t, dir, "tgt.env", "FOO=changed\n")

	results, err := MultiCompareAll(MultiCompareOptions{
		Baseline: base,
		Targets:  map[string]string{"staging": tgt},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !AnyTargetHasChanges(results) {
		t.Error("expected changes to be detected")
	}
}

func TestMultiCompareAll_MultipleTargets(t *testing.T) {
	dir := t.TempDir()
	base := writeEnvFile(t, dir, "base.env", "A=1\nB=2\n")
	same := writeEnvFile(t, dir, "same.env", "A=1\nB=2\n")
	diff := writeEnvFile(t, dir, "diff.env", "A=1\nB=99\n")

	results, err := MultiCompareAll(MultiCompareOptions{
		Baseline: base,
		Targets:  map[string]string{"same": same, "diff": diff},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !AnyTargetHasChanges(results) {
		t.Error("expected at least one target with changes")
	}
}

func TestMultiCompareAll_MissingBaseline(t *testing.T) {
	_, err := MultiCompareAll(MultiCompareOptions{
		Baseline: "/nonexistent/base.env",
		Targets:  map[string]string{"x": "/nonexistent/x.env"},
	})
	if err == nil {
		t.Error("expected error for missing baseline")
	}
}

func TestMultiCompareAll_NoTargets(t *testing.T) {
	_, err := MultiCompareAll(MultiCompareOptions{
		Baseline: "base.env",
		Targets:  map[string]string{},
	})
	if err == nil {
		t.Error("expected error when no targets provided")
	}
}
