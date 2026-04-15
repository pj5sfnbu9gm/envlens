package loader

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envlens-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", env["APP_ENV"])
	}
	if env["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %q", env["DEBUG"])
	}
}

func TestLoadFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
}

func TestLoadFile_StripQuotes(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/mydb"` + "\nSECRET='abc123'\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("unexpected DB_URL: %q", env["DB_URL"])
	}
	if env["SECRET"] != "abc123" {
		t.Errorf("unexpected SECRET: %q", env["SECRET"])
	}
}

func TestLoadFile_InvalidSyntax(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE\n")
	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for invalid syntax, got nil")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
