package masker_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/masker"
)

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"PRIVATE_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"LOG_LEVEL", false},
	}
	for _, tc := range cases {
		got := masker.IsSensitive(tc.key)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}

func TestMask_ShowsTrailingChars(t *testing.T) {
	opts := masker.MaskOptions{ShowChars: 4, Placeholder: "****"}
	got := masker.Mask("supersecretvalue", opts)
	want := "****alue"
	if got != want {
		t.Errorf("Mask() = %q, want %q", got, want)
	}
}

func TestMask_ShortValue(t *testing.T) {
	opts := masker.MaskOptions{ShowChars: 4, Placeholder: "****"}
	got := masker.Mask("abc", opts)
	if got != "****" {
		t.Errorf("expected full placeholder for short value, got %q", got)
	}
}

func TestMask_ZeroShowChars(t *testing.T) {
	opts := masker.MaskOptions{ShowChars: 0, Placeholder: "[REDACTED]"}
	got := masker.Mask("anyvalue", opts)
	if got != "[REDACTED]" {
		t.Errorf("expected full redaction, got %q", got)
	}
}

func TestMaskEnv(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abcd1234",
		"PORT":        "8080",
	}
	opts := masker.DefaultMaskOptions()
	masked := masker.MaskEnv(env, opts)

	if masked["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should not be masked")
	}
	if masked["PORT"] != "8080" {
		t.Errorf("PORT should not be masked")
	}
	if masked["DB_PASSWORD"] == "hunter2" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if masked["API_KEY"] == "abcd1234" {
		t.Errorf("API_KEY should be masked")
	}
}
