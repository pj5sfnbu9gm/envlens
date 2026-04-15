package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envlens/internal/snapshot"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	env := map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
	}

	if err := snapshot.Save(path, "prod", env); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if snap.Target != "prod" {
		t.Errorf("Target = %q, want %q", snap.Target, "prod")
	}

	if snap.Env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV = %q, want %q", snap.Env["APP_ENV"], "production")
	}

	if snap.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}

	if snap.Timestamp.Location() != time.UTC {
		t.Error("Timestamp should be UTC")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoad_MissingTarget(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notarget.json")

	data := `{"timestamp":"2024-01-01T00:00:00Z","env":{"KEY":"val"}}`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for missing target, got nil")
	}
}
