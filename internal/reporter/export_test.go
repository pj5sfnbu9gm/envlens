package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/reporter"
)

var exportSample = map[string]string{
	"APP_ENV": "staging",
	"PORT":    "8080",
	"DEBUG":   "true",
}

func TestReportExport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.ExportOptions{Format: "text", Target: "staging"}
	if err := reporter.ReportExport(&buf, exportSample, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "3 variable(s)") {
		t.Errorf("expected count in output, got:\n%s", out)
	}
	if !strings.Contains(out, "target: staging") {
		t.Errorf("expected target label, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV=staging") {
		t.Errorf("expected APP_ENV line, got:\n%s", out)
	}
}

func TestReportExport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.ExportOptions{Format: "json", Target: "production"}
	if err := reporter.ReportExport(&buf, exportSample, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["target"] != "production" {
		t.Errorf("expected target=production, got %v", result["target"])
	}
	if result["count"] != float64(3) {
		t.Errorf("expected count=3, got %v", result["count"])
	}
}

func TestReportExport_EmptyEnv(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.DefaultExportOptions()
	if err := reporter.ReportExport(&buf, map[string]string{}, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "0 variable(s)") {
		t.Errorf("expected zero count, got: %s", buf.String())
	}
}

func TestReportExport_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.ExportOptions{Format: "xml", Target: "staging"}
	err := reporter.ReportExport(&buf, exportSample, opts)
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("expected error to mention format name, got: %v", err)
	}
}
