package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

var sampleClusters = []differ.ClusterEntry{
	{
		Value:   "db.example.com",
		Keys:    []string{"CACHE_HOST", "DB_HOST"},
		Targets: []string{"prod", "staging"},
	},
}

func TestReportCluster_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	err := ReportCluster(&buf, sampleClusters, DefaultClusterReportOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "db.example.com") {
		t.Error("expected value in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key DB_HOST in output")
	}
	if !strings.Contains(out, "prod") {
		t.Error("expected target prod in output")
	}
}

func TestReportCluster_MaskValue(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultClusterReportOptions()
	opts.MaskVal = true
	_ = ReportCluster(&buf, sampleClusters, opts)
	if strings.Contains(buf.String(), "db.example.com") {
		t.Error("value should be masked")
	}
	if !strings.Contains(buf.String(), "***") {
		t.Error("expected mask placeholder")
	}
}

func TestReportCluster_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultClusterReportOptions()
	opts.Format = "json"
	err := ReportCluster(&buf, sampleClusters, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []struct {
		Value   string   `json:"value"`
		Keys    []string `json:"keys"`
		Targets []string `json:"targets"`
	}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0].Value != "db.example.com" {
		t.Errorf("unexpected value: %s", rows[0].Value)
	}
}

func TestReportCluster_Empty(t *testing.T) {
	var buf bytes.Buffer
	_ = ReportCluster(&buf, nil, DefaultClusterReportOptions())
	if !strings.Contains(buf.String(), "no clusters") {
		t.Error("expected no-clusters message")
	}
}
