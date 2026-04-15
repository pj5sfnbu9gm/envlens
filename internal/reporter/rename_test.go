package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/reporter"
)

func originalEnv() map[string]string {
	return map[string]string{
		"DB_HOST":    "localhost",
		"APP_SECRET": "s3cr3t",
		"LOG_LEVEL":  "info",
	}
}

func renamedEnv() map[string]string {
	return map[string]string{
		"DATABASE_HOST": "localhost",
		"APP_SECRET":    "s3cr3t",
		"LOG_LEVEL":     "info",
	}
}

func TestReportRename_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.DefaultRenameOptions(&buf)
	err := reporter.ReportRename(originalEnv(), renamedEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in text output")
	}
	if !strings.Contains(out, "DATABASE_HOST") {
		t.Error("expected DATABASE_HOST in text output")
	}
}

func TestReportRename_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.RenameOptions{Format: "json", Out: &buf}
	err := reporter.ReportRename(originalEnv(), renamedEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	renames, ok := result["renames"]
	if !ok {
		t.Fatal("expected 'renames' key in JSON output")
	}
	list, ok := renames.([]interface{})
	if !ok || len(list) == 0 {
		t.Fatal("expected non-empty renames array")
	}
}

func TestReportRename_NoChanges(t *testing.T) {
	env := originalEnv()
	var buf bytes.Buffer
	opts := reporter.DefaultRenameOptions(&buf)
	err := reporter.ReportRename(env, env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No keys renamed") {
		t.Errorf("expected 'No keys renamed' message, got: %s", buf.String())
	}
}
