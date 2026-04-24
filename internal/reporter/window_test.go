package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func sampleWindows() []differ.WindowResult {
	return []differ.WindowResult{
		{
			Labels: []string{"t0", "t1"},
			Results: []differ.Result{
				{Key: "HOST", Status: "changed", OldValue: "localhost", NewValue: "prod.example.com"},
				{Key: "NEW_KEY", Status: "added", NewValue: "yes"},
			},
		},
		{
			Labels:  []string{"t1", "t2"},
			Results: []differ.Result{},
		},
	}
}

func TestReportWindow_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultWindowReportOptions(&buf)
	if err := ReportWindow(sampleWindows(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "t0 → t1") {
		t.Error("expected window label 't0 → t1'")
	}
	if !strings.Contains(out, "~ HOST") {
		t.Error("expected changed HOST entry")
	}
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Error("expected added NEW_KEY entry")
	}
	if !strings.Contains(out, "(no changes)") {
		t.Error("expected '(no changes)' for empty window")
	}
}

func TestReportWindow_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultWindowReportOptions(&buf)
	opts.Format = "json"
	if err := ReportWindow(sampleWindows(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(parsed) != 2 {
		t.Errorf("expected 2 window entries, got %d", len(parsed))
	}
}

func TestReportWindow_Empty(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultWindowReportOptions(&buf)
	if err := ReportWindow(nil, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no windows") {
		t.Error("expected 'no windows' message for empty input")
	}
}
