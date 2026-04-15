package validator

import (
	"testing"
)

func TestValidate_NoFindings(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	findings := Validate(env, DefaultRules())
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": ""}
	findings := Validate(env, DefaultRules())
	if !containsRule(findings, "required-value") {
		t.Error("expected required-value finding for empty value")
	}
}

func TestValidate_InvalidKeyFormat(t *testing.T) {
	env := map[string]string{"my-key": "value"}
	findings := Validate(env, DefaultRules())
	if !containsRule(findings, "invalid-key-format") {
		t.Error("expected invalid-key-format finding")
	}
}

func TestValidate_LeadingTrailingSpace(t *testing.T) {
	env := map[string]string{"APP_NAME": "  myapp "}
	findings := Validate(env, DefaultRules())
	if !containsRule(findings, "leading-trailing-space") {
		t.Error("expected leading-trailing-space finding")
	}
}

func TestValidate_CustomRule(t *testing.T) {
	customRule := Rule{
		Name:    "no-localhost",
		Message: "value should not be localhost in production",
		Check:   func(_, v string) bool { return v == "localhost" },
	}
	env := map[string]string{"DB_HOST": "localhost"}
	findings := Validate(env, []Rule{customRule})
	if len(findings) != 1 || findings[0].Rule != "no-localhost" {
		t.Errorf("expected no-localhost finding, got %+v", findings)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	// lowercase key AND empty value — should produce two findings
	env := map[string]string{"bad_key": ""}
	findings := Validate(env, DefaultRules())
	if len(findings) < 2 {
		t.Errorf("expected at least 2 findings, got %d", len(findings))
	}
}

func TestValidate_EmptyEnv(t *testing.T) {
	findings := Validate(map[string]string{}, DefaultRules())
	if len(findings) != 0 {
		t.Errorf("expected no findings for empty env, got %d", len(findings))
	}
}

// containsRule is a helper that checks whether any finding has the given rule name.
func containsRule(findings []Finding, rule string) bool {
	for _, f := range findings {
		if f.Rule == rule {
			return true
		}
	}
	return false
}
