package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/envlens/internal/differ"
)

func sampleSlices() []differ.SliceResult {
	return []differ.SliceResult{
		{
			Prefix: "DB",
			Results: []differ.Result{
				{Key: "DB_HOST", Status: differ.StatusChanged, OldValue: "localhost", NewValue: "prod-db"},
			},
		},
		{
			Prefix: "APP",
			Results: []differ.Result{
				{Key: "APP_ENV", Status: differ.StatusAdded, NewValue: "production"},
				{Key: "APP_DEBUG", Status: differ.StatusRemoved, OldValue: "true"},
			},
		},
	}
}

func TestReportSlice_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultSliceReportOptions()
	opts.Writer = &buf
	if err := ReportSlice(sampleSlices(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[DB]") {
		t.Error("expected [DB] section header")
	}
	if !strings.Contains(out, "~ DB_HOST") {
		t.Error("expected changed marker for DB_HOST")
	}
	if !strings.Contains(out, "+ APP_ENV") {
		t.Error("expected added marker for APP_ENV")
	}
	if !strings.Contains(out, "- APP_DEBUG") {
		t.Error("expected removed marker for APP_DEBUG")
	}
}

func TestReportSlice_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultSliceReportOptions()
	opts.Format = "json"
	opts.Writer = &buf
	if err := ReportSlice(sampleSlices(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []struct {
		Prefix  string           `json:"prefix"`
		Results []differ.Result  `json:"results"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 groups, got %d", len(out))
	}
	if out[0].Prefix != "DB" {
		t.Errorf("expected first prefix DB, got %s", out[0].Prefix)
	}
}

func TestReportSlice_Empty(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultSliceReportOptions()
	opts.Writer = &buf
	if err := ReportSlice(nil, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no slice changes") {
		t.Error("expected 'no slice changes' message")
	}
}
