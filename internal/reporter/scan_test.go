package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/reporter"
	"github.com/user/envlens/internal/scanner"
)

func sampleFindings() []scanner.Finding {
	return []scanner.Finding{
		{Key: "API_KEY", Value: "abc", Severity: "warning", Message: "key \"API_KEY\" appears sensitive but has a suspiciously short value"},
		{Key: "DB_PASSWORD", Value: "secret", Severity: "error", Message: "key \"DB_PASSWORD\" has a well-known default value"},
	}
}

func TestReportScan_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.DefaultScanOptions()
	opts.Out = &buf
	if err := reporter.ReportScan(sampleFindings(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[error]") {
		t.Error("expected [error] in output")
	}
	if !strings.Contains(out, "[warning]") {
		t.Error("expected [warning] in output")
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in output")
	}
}

func TestReportScan_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.DefaultScanOptions()
	opts.Format = "json"
	opts.Out = &buf
	if err := reporter.ReportScan(sampleFindings(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 findings, got %d", len(result))
	}
}

func TestReportScan_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.DefaultScanOptions()
	opts.Out = &buf
	if err := reporter.ReportScan(nil, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no issues") {
		t.Error("expected 'no issues' message for empty findings")
	}
}
