package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func sampleChain() []differ.ChainResult {
	return []differ.ChainResult{
		{
			From: "dev",
			To:   "staging",
			Results: []differ.Result{
				{Key: "HOST", Status: differ.StatusChanged, OldValue: "localhost", NewValue: "staging.host"},
				{Key: "PORT", Status: differ.StatusUnchanged, OldValue: "8080", NewValue: "8080"},
			},
		},
	}
}

func TestReportChain_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultChainOptions()
	opts.Writer = &buf
	err := ReportChain(sampleChain(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "dev -> staging") {
		t.Error("expected step header")
	}
	if !strings.Contains(out, "changed") {
		t.Error("expected changed status")
	}
	if strings.Contains(out, "unchanged") {
		t.Error("should hide unchanged by default")
	}
}

func TestReportChain_ShowUnchanged(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultChainOptions()
	opts.Writer = &buf
	opts.ShowUnchanged = true
	ReportChain(sampleChain(), opts)
	out := buf.String()
	if !strings.Contains(out, "unchanged") {
		t.Error("expected unchanged entries when ShowUnchanged=true")
	}
}

func TestReportChain_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultChainOptions()
	opts.Writer = &buf
	opts.Format = "json"
	ReportChain(sampleChain(), opts)
	if !strings.Contains(buf.String(), "\"From\"") {
		t.Error("expected JSON output with From field")
	}
}

func TestReportChain_Empty(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultChainOptions()
	opts.Writer = &buf
	ReportChain(nil, opts)
	if !strings.Contains(buf.String(), "no chain steps") {
		t.Error("expected no-steps message")
	}
}
