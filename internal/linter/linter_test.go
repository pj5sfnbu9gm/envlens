package linter

import (
	"testing"
)

func TestLint_NoFindings(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	findings := Lint(env, DefaultRules())
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": ""}
	findings := Lint(env, []Rule{RuleNoEmptyValue})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != SeverityWarning {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	findings := Lint(env, []Rule{RuleUppercaseKey})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "db_host" {
		t.Errorf("unexpected key: %s", findings[0].Key)
	}
}

func TestLint_WhitespaceInKey(t *testing.T) {
	env := map[string]string{"BAD KEY": "value"}
	findings := Lint(env, []Rule{RuleNoWhitespaceInKey})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != SeverityError {
		t.Errorf("expected error severity, got %s", findings[0].Severity)
	}
}

func TestLint_LeadingUnderscore(t *testing.T) {
	env := map[string]string{"_INTERNAL": "secret"}
	findings := Lint(env, []Rule{RuleNoLeadingUnderscore})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != SeverityInfo {
		t.Errorf("expected info severity, got %s", findings[0].Severity)
	}
}

func TestLint_CustomRule(t *testing.T) {
	noDebug := func(key, value string) *Finding {
		if key == "DEBUG" && value == "true" {
			return &Finding{Key: key, Message: "DEBUG should not be true in production", Severity: SeverityError}
		}
		return nil
	}
	env := map[string]string{"DEBUG": "true", "APP_ENV": "production"}
	findings := Lint(env, []Rule{noDebug})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Message != "DEBUG should not be true in production" {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestLint_MultipleViolations(t *testing.T) {
	env := map[string]string{
		"bad key":   "",
		"_internal": "val",
	}
	findings := Lint(env, DefaultRules())
	// each key will hit multiple rules
	if len(findings) == 0 {
		t.Error("expected findings but got none")
	}
}
