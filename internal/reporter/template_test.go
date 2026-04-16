package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestReportTemplate_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	res := TemplateResult{Output: "app=envlens port=8080", MissingKeys: nil}
	err := ReportTemplate(res, DefaultTemplateOptions(&buf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "app=envlens port=8080") {
		t.Errorf("output missing rendered text: %q", buf.String())
	}
}

func TestReportTemplate_TextFormat_MissingKeys(t *testing.T) {
	var buf bytes.Buffer
	res := TemplateResult{Output: "", MissingKeys: []string{"SECRET", "API_KEY"}}
	err := ReportTemplate(res, DefaultTemplateOptions(&buf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "missing keys") {
		t.Errorf("expected missing keys section: %q", out)
	}
	if !strings.Contains(out, "API_KEY") || !strings.Contains(out, "SECRET") {
		t.Errorf("expected both missing keys listed: %q", out)
	}
}

func TestReportTemplate_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	res := TemplateResult{Output: "hello", MissingKeys: []string{"X"}}
	opts := DefaultTemplateOptions(&buf)
	opts.Format = "json"
	err := ReportTemplate(res, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["output"] != "hello" {
		t.Errorf("unexpected output field: %v", parsed["output"])
	}
	keys, _ := parsed["missing_keys"].([]interface{})
	if len(keys) != 1 || keys[0] != "X" {
		t.Errorf("unexpected missing_keys: %v", keys)
	}
}

func TestReportTemplate_NoMissingKeys_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	res := TemplateResult{Output: "clean output", MissingKeys: []string{}}
	err := ReportTemplate(res, DefaultTemplateOptions(&buf))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "missing keys") {
		t.Errorf("should not show missing keys section when empty")
	}
}
