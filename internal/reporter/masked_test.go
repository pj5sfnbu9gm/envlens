package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/masker"
	"github.com/yourorg/envlens/internal/reporter"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "envlens",
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "key-abcd-1234",
		"PORT":        "9000",
	}
}

func TestReportMasked_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.DefaultMaskedOptions()
	if err := reporter.ReportMasked(&buf, sampleEnv(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "APP_NAME=envlens") {
		t.Errorf("expected APP_NAME unmasked, got:\n%s", out)
	}
	if strings.Contains(out, "supersecret") {
		t.Errorf("DB_PASSWORD value should be masked, got:\n%s", out)
	}
	if strings.Contains(out, "key-abcd-1234") {
		t.Errorf("API_KEY value should be masked, got:\n%s", out)
	}
}

func TestReportMasked_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.MaskedOptions{
		Format: "json",
		Mask:   masker.DefaultMaskOptions(),
	}
	if err := reporter.ReportMasked(&buf, sampleEnv(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result["APP_NAME"] != "envlens" {
		t.Errorf("APP_NAME should be unmasked in JSON")
	}
	if result["DB_PASSWORD"] == "supersecret" {
		t.Errorf("DB_PASSWORD should be masked in JSON")
	}
}

func TestReportMasked_FullRedaction(t *testing.T) {
	var buf bytes.Buffer
	opts := reporter.MaskedOptions{
		Format: "text",
		Mask:   masker.MaskOptions{ShowChars: 0, Placeholder: "[HIDDEN]"},
	}
	if err := reporter.ReportMasked(&buf, sampleEnv(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "supersecret") || strings.Contains(out, "key-abcd") {
		t.Errorf("expected full redaction, got:\n%s", out)
	}
	if !strings.Contains(out, "[HIDDEN]") {
		t.Errorf("expected [HIDDEN] placeholder, got:\n%s", out)
	}
}
