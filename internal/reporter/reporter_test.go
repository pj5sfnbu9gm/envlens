package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func makeResults() []differ.Result {
	return []differ.Result{
		{Key: "APP_ENV", Status: differ.StatusUnchanged, FromValue: "production", ToValue: "production"},
		{Key: "DB_HOST", Status: differ.StatusChanged, FromValue: "db-old", ToValue: "db-new"},
		{Key: "NEW_KEY", Status: differ.StatusAdded, FromValue: "", ToValue: "value1"},
		{Key: "OLD_KEY", Status: differ.StatusRemoved, FromValue: "gone", ToValue: ""},
	}
}

func TestReport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := Options{Format: FormatText, NoColor: true, Writer: &buf}
	err := Report(makeResults(), "staging", "production", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "staging") || !strings.Contains(out, "production") {
		t.Error("expected header to contain source and target names")
	}
	if !strings.Contains(out, "+ NEW_KEY=value1") {
		t.Error("expected added key line")
	}
	if !strings.Contains(out, "- OLD_KEY=gone") {
		t.Error("expected removed key line")
	}
	if !strings.Contains(out, "~ DB_HOST") {
		t.Error("expected changed key line")
	}
	if strings.Contains(out, "APP_ENV") {
		t.Error("unchanged keys should not appear in text output")
	}
}

func TestReport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := Options{Format: FormatJSON, Writer: &buf}
	err := Report(makeResults(), "staging", "production", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var report jsonReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}
	if report.From != "staging" || report.To != "production" {
		t.Errorf("unexpected from/to: %s / %s", report.From, report.To)
	}
	if report.Summary.Added != 1 || report.Summary.Removed != 1 || report.Summary.Changed != 1 {
		t.Errorf("unexpected summary: %+v", report.Summary)
	}
	if len(report.Changes) != 4 {
		t.Errorf("expected 4 changes, got %d", len(report.Changes))
	}
}

func TestReport_NoChanges(t *testing.T) {
	results := []differ.Result{
		{Key: "FOO", Status: differ.StatusUnchanged, FromValue: "bar", ToValue: "bar"},
	}
	var buf bytes.Buffer
	opts := Options{Format: FormatText, NoColor: true, Writer: &buf}
	err := Report(results, "a", "b", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Error("expected no-differences message")
	}
}
