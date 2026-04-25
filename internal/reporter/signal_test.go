package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

var sampleSignalEntries = []differ.SignalEntry{
	{Key: "DB_HOST", ChangeCount: 3, Targets: []string{"dev", "prod", "staging"}},
	{Key: "API_KEY", ChangeCount: 2, Targets: []string{"dev", "prod"}},
}

func TestReportSignal_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultSignalReportOptions(&buf)
	if err := ReportSignal(sampleSignalEntries, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "changes=3") {
		t.Error("expected changes=3 in output")
	}
	if !strings.Contains(out, "dev,prod,staging") {
		t.Error("expected target list in output")
	}
}

func TestReportSignal_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultSignalReportOptions(&buf)
	opts.Format = "json"
	if err := ReportSignal(sampleSignalEntries, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded []differ.SignalEntry
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(decoded) != 2 {
		t.Errorf("expected 2 entries, got %d", len(decoded))
	}
	if decoded[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", decoded[0].Key)
	}
}

func TestReportSignal_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultSignalReportOptions(&buf)
	if err := ReportSignal(nil, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no high-signal") {
		t.Error("expected empty message")
	}
}
