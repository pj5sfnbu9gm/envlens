package auditor_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/auditor"
)

func TestAudit_NoFindings(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
	}
	findings := auditor.Audit(env, auditor.DefaultRules())
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestAudit_EmptyValue(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "",
	}
	findings := auditor.Audit(env, auditor.DefaultRules())
	if !containsRule(findings, "no-empty-value") {
		t.Error("expected no-empty-value finding")
	}
}

func TestAudit_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"app_secret": "abc123",
	}
	findings := auditor.Audit(env, auditor.DefaultRules())
	if !containsRule(findings, "uppercase-key") {
		t.Error("expected uppercase-key finding")
	}
}

func TestAudit_WhitespaceInKey(t *testing.T) {
	env := map[string]string{
		"BAD KEY": "value",
	}
	findings := auditor.Audit(env, auditor.DefaultRules())
	if !containsRule(findings, "no-whitespace-key") {
		t.Error("expected no-whitespace-key finding")
	}
}

func TestAudit_CustomRule(t *testing.T) {
	env := map[string]string{
		"SECRET": "hunter2",
	}
	rules := []auditor.Rule{
		{
			Name: "min-length",
			Check: func(key, value string) error {
				if len(value) < 10 {
					return fmt.Errorf("key %q value too short", key)
				}
				return nil
			},
		},
	}
	findings := auditor.Audit(env, rules)
	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}
}

func TestAudit_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "",
		"A_KEY": "",
		"M_KEY": "",
	}
	findings := auditor.Audit(env, auditor.DefaultRules())
	if len(findings) < 3 {
		t.Fatal("expected at least 3 findings")
	}
	if findings[0].Key > findings[1].Key || findings[1].Key > findings[2].Key {
		t.Error("findings are not sorted by key")
	}
}

func containsRule(findings []auditor.Finding, rule string) bool {
	for _, f := range findings {
		if f.Rule == rule {
			return true
		}
	}
	return false
}
