package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlens/internal/resolver"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestResolveTargets_Basic(t *testing.T) {
	dir := t.TempDir()
	specs := []string{"prod=prod.env", "staging=staging.env"}

	targets, err := resolver.ResolveTargets(specs, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
	if targets[0].Name != "prod" {
		t.Errorf("expected name prod, got %q", targets[0].Name)
	}
	if targets[1].Name != "staging" {
		t.Errorf("expected name staging, got %q", targets[1].Name)
	}
}

func TestResolveTargets_InvalidSpec(t *testing.T) {
	cases := []string{"noequalssign", "=noname", "nopath="}
	for _, spec := range cases {
		_, err := resolver.ResolveTargets([]string{spec}, "/tmp")
		if err == nil {
			t.Errorf("expected error for spec %q", spec)
		}
	}
}

func TestResolveTargets_DuplicateName(t *testing.T) {
	_, err := resolver.ResolveTargets([]string{"prod=a.env", "prod=b.env"}, "/tmp")
	if err == nil {
		t.Error("expected error for duplicate target name")
	}
}

func TestResolveTargets_EmptySpecs(t *testing.T) {
	_, err := resolver.ResolveTargets([]string{}, "/tmp")
	if err == nil {
		t.Error("expected error for empty specs")
	}
}

func TestValidatePaths_OK(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "prod.env", "KEY=val\n")

	targets := []resolver.Target{{Name: "prod", Path: filepath.Join(dir, "prod.env")}}
	if err := resolver.ValidatePaths(targets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidatePaths_Missing(t *testing.T) {
	targets := []resolver.Target{{Name: "prod", Path: "/nonexistent/prod.env"}}
	if err := resolver.ValidatePaths(targets); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidatePaths_Directory(t *testing.T) {
	dir := t.TempDir()
	targets := []resolver.Target{{Name: "prod", Path: dir}}
	if err := resolver.ValidatePaths(targets); err == nil {
		t.Error("expected error when path is a directory")
	}
}
