package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/freezer"
)

func frozenSample(t *testing.T) *freezer.FrozenEnv {
	t.Helper()
	f, err := freezer.Freeze(map[string]string{
		"HOST": "localhost",
		"PORT": "9000",
	}, freezer.DefaultOptions())
	if err != nil {
		t.Fatalf("freeze: %v", err)
	}
	return f
}

func TestReportFreeze_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	f := frozenSample(t)
	err := ReportFreeze(f, DefaultFreezeOptions(&buf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Frozen environment") {
		t.Error("expected header in text output")
	}
	if !strings.Contains(out, "HOST=localhost") {
		t.Error("expected HOST entry")
	}
}

func TestReportFreeze_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	f := frozenSample(t)
	opts := FreezeOptions{Format: "json", Writer: &buf}
	err := ReportFreeze(f, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out struct {
		Total int               `json:"total"`
		Env   map[string]string `json:"env"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Total != 2 {
		t.Errorf("expected total=2, got %d", out.Total)
	}
	if out.Env["PORT"] != "9000" {
		t.Errorf("expected PORT=9000")
	}
}

func TestReportFreeze_EmptyEnv(t *testing.T) {
	var buf bytes.Buffer
	f, _ := freezer.Freeze(map[string]string{}, freezer.DefaultOptions())
	err := ReportFreeze(f, DefaultFreezeOptions(&buf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "0 keys") {
		t.Error("expected '0 keys' in output")
	}
}
