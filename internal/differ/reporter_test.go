package differ

import (
	"bytes"
	"strings"
	"testing"
)

func sampleResults() []Result {
	return []Result{
		{Key: "HOST", Status: StatusUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "PORT", Status: StatusChanged, OldValue: "8080", NewValue: "9090"},
		{Key: "NEW_KEY", Status: StatusAdded, OldValue: "", NewValue: "hello"},
		{Key: "OLD_KEY", Status: StatusRemoved, OldValue: "bye", NewValue: ""},
	}
}

func TestReportDiff_TextAll(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions()
	if err := ReportDiff(&buf, sampleResults(), opts); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY=hello") {
		t.Errorf("expected added line, got:\n%s", out)
	}
	if !strings.Contains(out, "- OLD_KEY=bye") {
		t.Errorf("expected removed line, got:\n%s", out)
	}
	if !strings.Contains(out, "~ PORT: 8080 -> 9090") {
		t.Errorf("expected changed line, got:\n%s", out)
	}
}

func TestReportDiff_TextFilterChanged(t *testing.T) {
	var buf bytes.Buffer
	opts := ReportOptions{Format: "text", ShowOnly: "changed"}
	ReportDiff(&buf, sampleResults(), opts)
	out := buf.String()
	if strings.Contains(out, "NEW_KEY") || strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected only changed keys, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got:\n%s", out)
	}
}

func TestReportDiff_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := ReportOptions{Format: "json", ShowOnly: "all"}
	if err := ReportDiff(&buf, sampleResults(), opts); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"Key\"") {
		t.Errorf("expected JSON output, got:\n%s", out)
	}
}

func TestReportDiff_NoResults(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultReportOptions()
	ReportDiff(&buf, []Result{}, opts)
	if !strings.Contains(buf.String(), "no differences found") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestReportDiff_FilterAdded(t *testing.T) {
	var buf bytes.Buffer
	opts := ReportOptions{Format: "text", ShowOnly: "added"}
	ReportDiff(&buf, sampleResults(), opts)
	out := buf.String()
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY, got: %s", out)
	}
	if strings.Contains(out, "PORT") {
		t.Errorf("unexpected PORT in added-only output")
	}
}
