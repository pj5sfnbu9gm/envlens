package scanner_test

import (
	"testing"

	"github.com/user/envlens/internal/scanner"
)

func TestScan_NoFindings(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	findings := scanner.Scan(env, scanner.Options{})
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestScan_ShortSensitiveValue(t *testing.T) {
	env := map[string]string{"API_KEY": "abc"}
	findings := scanner.Scan(env, scanner.Options{})
	if len(findings) == 0 {
		t.Fatal("expected finding for short sensitive value")
	}
	if findings[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", findings[0].Severity)
	}
}

func TestScan_LocalhostURL(t *testing.T) {
	env := map[string]string{"DB_HOST": "http://localhost:5432"}
	findings := scanner.Scan(env, scanner.Options{})
	if len(findings) == 0 {
		t.Fatal("expected finding for localhost URL")
	}
	if findings[0].Severity != "info" {
		t.Errorf("expected info, got %s", findings[0].Severity)
	}
}

func TestScan_DefaultPassword(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "password"}
	findings := scanner.Scan(env, scanner.Options{})
	found := false
	for _, f := range findings {
		if f.Severity == "error" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected error severity finding for default password")
	}
}

func TestScan_CustomRule(t *testing.T) {
	called := false
	rule := func(key, value string) []scanner.Finding {
		called = true
		return nil
	}
	env := map[string]string{"FOO": "bar"}
	scanner.Scan(env, scanner.Options{Rules: []scanner.Rule{rule}})
	if !called {
		t.Fatal("expected custom rule to be called")
	}
}

func TestScan_MultipleFindings(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"API_URL":     "http://127.0.0.1/api",
	}
	findings := scanner.Scan(env, scanner.Options{})
	if len(findings) < 2 {
		t.Fatalf("expected at least 2 findings, got %d", len(findings))
	}
}
