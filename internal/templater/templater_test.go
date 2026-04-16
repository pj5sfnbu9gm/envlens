package templater

import (
	"strings"
	"testing"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "envlens",
		"PORT":     "8080",
	}
}

func TestRender_BasicSubstitution(t *testing.T) {
	result, err := Render(`app={{ env "APP_NAME" }} port={{ env "PORT" }}`, sampleEnv(), DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "app=envlens port=8080" {
		t.Errorf("got %q", result.Output)
	}
}

func TestRender_EnvOr_Present(t *testing.T) {
	result, err := Render(`{{ envOr "APP_NAME" "default" }}`, sampleEnv(), DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "envlens" {
		t.Errorf("got %q", result.Output)
	}
}

func TestRender_EnvOr_Missing(t *testing.T) {
	result, err := Render(`{{ envOr "MISSING" "fallback" }}`, sampleEnv(), DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "fallback" {
		t.Errorf("got %q", result.Output)
	}
}

func TestRender_MissingKey_Lenient(t *testing.T) {
	opts := DefaultOptions()
	opts.FailOnMissing = false
	result, err := Render(`{{ env "GHOST" }}`, sampleEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.MissingKeys) != 1 || result.MissingKeys[0] != "GHOST" {
		t.Errorf("expected missing key GHOST, got %v", result.MissingKeys)
	}
}

func TestRender_MissingKey_Strict(t *testing.T) {
	opts := DefaultOptions()
	opts.FailOnMissing = true
	_, err := Render(`{{ env "GHOST" }}`, sampleEnv(), opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRender_UpperLower(t *testing.T) {
	result, err := Render(`{{ upper "hello" }}-{{ lower "WORLD" }}`, sampleEnv(), DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "HELLO-world" {
		t.Errorf("got %q", result.Output)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	_, err := Render(`{{ unclosed`, sampleEnv(), DefaultOptions())
	if err == nil {
		t.Fatal("expected parse error")
	}
	if !strings.Contains(err.Error(), "parse template") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRender_CustomDelimiters(t *testing.T) {
	opts := DefaultOptions()
	opts.LeftDelim = "[["
	opts.RightDelim = "]]"
	result, err := Render(`[[ env "APP_NAME" ]]`, sampleEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "envlens" {
		t.Errorf("got %q", result.Output)
	}
}

func TestRender_EmptyEnv(t *testing.T) {
	opts := DefaultOptions()
	opts.FailOnMissing = false
	result, err := Render(`{{ env "APP_NAME" }}`, map[string]string{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.MissingKeys) != 1 || result.MissingKeys[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME in missing keys, got %v", result.MissingKeys)
	}
}
