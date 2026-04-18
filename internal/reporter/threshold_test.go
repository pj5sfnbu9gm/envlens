package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func sampleThresholdResults() []differ.ThresholdResult {
	return []differ.ThresholdResult{
		{
			Target: "prod",
			Count:  2,
			Results: []differ.Result{
				{Key: "FOO", Status: "changed", OldValue: "a", NewValue: "b"},
				{Key: "BAR", Status: "added", NewValue: "new"},
			},
		},
	}
}

func TestReportThreshold_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultThresholdReportOptions(&buf)
	if err := ReportThreshold(sampleThresholdResults(), opts); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Error("expected target name")
	}
	if !strings.Contains(out, "FOO") {
		t.Error("expected key FOO")
	}
	if !strings.Contains(out, "BAR") {
		t.Error("expected key BAR")
	}
}

func TestReportThreshold_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultThresholdReportOptions(&buf)
	opts.Format = "json"
	if err := ReportThreshold(sampleThresholdResults(), opts); err != nil {
		t.Fatal(err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}

func TestReportThreshold_NoResults(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultThresholdReportOptions(&buf)
	if err := ReportThreshold(nil, opts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no targets") {
		t.Error("expected no-targets message")
	}
}
