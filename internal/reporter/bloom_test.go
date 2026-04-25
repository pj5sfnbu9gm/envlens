package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

var sampleBloom = []differ.BloomEntry{
	{Key: "API_KEY", PresentIn: []string{"dev", "prod"}, AbsentIn: nil},
	{Key: "DEBUG", PresentIn: []string{"dev"}, AbsentIn: []string{"prod"}},
	{Key: "FEATURE_X", PresentIn: []string{"prod"}, AbsentIn: []string{"dev"}},
}

func TestReportBloom_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultBloomReportOptions()
	if err := ReportBloom(&buf, sampleBloom, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
	if !strings.Contains(out, "absent") {
		t.Error("expected 'absent' label for keys missing in some targets")
	}
}

func TestReportBloom_ShowGapsOnly(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultBloomReportOptions()
	opts.ShowGapsOnly = true
	if err := ReportBloom(&buf, sampleBloom, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "API_KEY") {
		t.Error("API_KEY should be filtered out (no gaps)")
	}
	if !strings.Contains(out, "DEBUG") {
		t.Error("DEBUG should appear (has gaps)")
	}
}

func TestReportBloom_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultBloomReportOptions()
	opts.Format = "json"
	if err := ReportBloom(&buf, sampleBloom, opts); err != nil {
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
		if _, ok := r["absent_in"]; !ok {
			t.Error("expected absent_in field in JSON")
		}
	}
}

func TestReportBloom_Empty(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultBloomReportOptions()
	if err := ReportBloom(&buf, nil, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no bloom entries") {
		t.Error("expected empty message")
	}
}
