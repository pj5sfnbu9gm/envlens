package encoder

import (
	"encoding/base64"
	"strings"
	"testing"
)

var sample = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"LOG_LEVEL": "info",
}

func TestEncode_JSON(t *testing.T) {
	opts := DefaultOptions()
	out, err := Encode(sample, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"APP_ENV"`) {
		t.Errorf("expected APP_ENV in JSON output, got: %s", out)
	}
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestEncode_Base64(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatBase64
	out, err := Encode(sample, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		t.Fatalf("base64 decode failed: %v", err)
	}
	if !strings.Contains(string(decoded), "APP_ENV") {
		t.Errorf("decoded base64 missing APP_ENV: %s", decoded)
	}
}

func TestEncode_CSV(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatCSV
	out, err := Encode(sample, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(out, "\n")
	if lines[0] != "key,value" {
		t.Errorf("expected header 'key,value', got: %s", lines[0])
	}
	if len(lines) != 4 { // header + 3 keys
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestEncode_SortedKeys(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatCSV
	opts.SortKeys = true
	out, err := Encode(sample, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(out, "\n")
	if !strings.HasPrefix(lines[1], "APP_ENV") {
		t.Errorf("expected APP_ENV first after header, got: %s", lines[1])
	}
}

func TestEncode_UnknownFormat(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = "xml"
	_, err := Encode(sample, opts)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestEncode_EmptyEnv(t *testing.T) {
	opts := DefaultOptions()
	out, err := Encode(map[string]string{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "{}" {
		t.Errorf("expected '{}', got: %s", out)
	}
}
