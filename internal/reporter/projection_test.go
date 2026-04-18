package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func projectionSample() []differ.DiffResult {
	return []differ.DiffResult{
		{Key: "APP_PORT", Status: differ.StatusChanged, OldValue: "8080", NewValue: "9090"},
		{Key: "DB_HOST", Status: differ.StatusAdded, OldValue: "", NewValue: "db.prod"},
		{Key: "LOG_LEVEL", Status: differ.StatusUnchanged, OldValue: "info", NewValue: "info"},
	}
}

func TestReportProjection_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	err := ReportProjection(projectionSample(), DefaultProjectionReportOptions(&buf))
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "~ APP_PORT") {
		t.Errorf("expected changed marker, got: %s", out)
	}
	if !strings.Contains(out, "+ DB_HOST") {
		t.Errorf("expected added marker, got: %s", out)
	}
	if !strings.Contains(out, "  LOG_LEVEL") {
		t.Errorf("expected unchanged marker, got: %s", out)
	}
}

func TestReportProjection_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := ReportProjection(projectionSample(), ProjectionReportOptions{Format: "json", Writer: &buf})
	if err != nil {
		t.Fatal(err)
	}
	var out []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestReportProjection_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := ReportProjection(nil, DefaultProjectionReportOptions(&buf))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no projected results") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
