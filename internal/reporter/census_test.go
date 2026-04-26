package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/differ"
)

func sampleCensusEntries() []differ.CensusEntry {
	return []differ.CensusEntry{
		{Key: "DB_HOST", Count: 3, Targets: []string{"dev", "prod", "staging"}, Coverage: 1.0},
		{Key: "API_KEY", Count: 2, Targets: []string{"prod", "staging"}, Coverage: 0.667},
		{Key: "DEBUG", Count: 1, Targets: []string{"dev"}, Coverage: 0.333},
	}
}

func TestReportCensus_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultCensusReportOptions(3)
	opts.Out = &buf

	if err := ReportCensus(sampleCensusEntries(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "GAP") {
		t.Error("expected GAP marker for partial keys")
	}
	if strings.Count(out, "GAP") != 2 {
		t.Errorf("expected 2 GAP markers, got %d", strings.Count(out, "GAP"))
	}
}

func TestReportCensus_ShowGapsOnly(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultCensusReportOptions(3)
	opts.Out = &buf
	opts.ShowGapsOnly = true

	if err := ReportCensus(sampleCensusEntries(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "DB_HOST") {
		t.Error("DB_HOST should be excluded (universal key)")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("API_KEY should appear")
	}
}

func TestReportCensus_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultCensusReportOptions(3)
	opts.Out = &buf
	opts.Format = "json"

	if err := ReportCensus(sampleCensusEntries(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(rows))
	}
	for _, r := range rows {
		if _, ok := r["gap"]; !ok {
			t.Error("expected 'gap' field in JSON output")
		}
	}
}

func TestReportCensus_Empty(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultCensusReportOptions(3)
	opts.Out = &buf

	if err := ReportCensus(nil, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no census data") {
		t.Error("expected empty message")
	}
}
