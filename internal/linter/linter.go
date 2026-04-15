package linter

import (
	"fmt"
	"strings"
)

// Severity represents the severity level of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint finding for an environment variable.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Rule is a function that inspects a key-value pair and returns a Finding if
// the rule is violated, or nil if the pair is compliant.
type Rule func(key, value string) *Finding

// DefaultRules returns the built-in set of lint rules.
func DefaultRules() []Rule {
	return []Rule{
		RuleNoEmptyValue,
		RuleUppercaseKey,
		RuleNoWhitespaceInKey,
		RuleNoLeadingUnderscore,
	}
}

// Lint runs all provided rules against the given env map and returns all findings.
func Lint(env map[string]string, rules []Rule) []Finding {
	findings := []Finding{}
	for k, v := range env {
		for _, rule := range rules {
			if f := rule(k, v); f != nil {
				findings = append(findings, *f)
			}
		}
	}
	return findings
}

// RuleNoEmptyValue flags keys with empty values.
func RuleNoEmptyValue(key, value string) *Finding {
	if strings.TrimSpace(value) == "" {
		return &Finding{Key: key, Message: "value is empty", Severity: SeverityWarning}
	}
	return nil
}

// RuleUppercaseKey flags keys that contain lowercase letters.
func RuleUppercaseKey(key, _ string) *Finding {
	if key != strings.ToUpper(key) {
		return &Finding{Key: key, Message: "key should be uppercase", Severity: SeverityWarning}
	}
	return nil
}

// RuleNoWhitespaceInKey flags keys that contain whitespace characters.
func RuleNoWhitespaceInKey(key, _ string) *Finding {
	if strings.ContainsAny(key, " \t") {
		return &Finding{Key: key, Message: "key contains whitespace", Severity: SeverityError}
	}
	return nil
}

// RuleNoLeadingUnderscore flags keys that start with an underscore.
func RuleNoLeadingUnderscore(key, _ string) *Finding {
	if strings.HasPrefix(key, "_") {
		return &Finding{
			Key:      key,
			Message:  fmt.Sprintf("key %q starts with underscore", key),
			Severity: SeverityInfo,
		}
	}
	return nil
}
