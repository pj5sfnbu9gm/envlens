package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func sampleLabelEntries() []differ.LabelEntry {
	return []differ.LabelEntry{
		{Key: "APP_ENV", Status: "changed", Label: "environment", Old: "staging", New: "production"},
		{Key: "NEW_KEY", Status: "added", Label: "unlabeled", Old: "", New: "val"},
	}
}

func TestReportLabel_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultLabelReportOptions()
	opts.Writer = &buf
	err := ReportLabel(sampleLabelEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in output")
	}
	if !strings.Contains(out, "environment") {
		t.Error("expected label 'environment' in output")
	}
	if !strings.Contains(out, "KEY") {
		t.Error("expected header row in output")
	}
}

func TestReportLabel_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultLabelReportOptions()
	opts.Format = "json"
	opts.Writer = &buf
	err := ReportLabel(sampleLabelEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}
	if rows[0]["label"] != "environment" {
		t.Errorf("expected label 'environment', got %q", rows[0]["label"])
	}
}

func TestReportLabel_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultLabelReportOptions()
	opts.Writer = &buf
	err := ReportLabel([]differ.LabelEntry{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no labeled changes") {
		t.Error("expected 'no labeled changes' message")
	}
}
