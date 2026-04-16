package scanner

import (
	"fmt"
	"strings"
)

// Finding represents a single scan result for a key/value pair.
type Finding struct {
	Key      string
	Value    string
	Severity string // "error", "warning", "info"
	Message  string
}

// Rule is a function that inspects a key/value pair and returns findings.
type Rule func(key, value string) []Finding

// Options controls scanner behaviour.
type Options struct {
	Rules []Rule
}

// DefaultRules returns the built-in scan rules.
func DefaultRules() []Rule {
	return []Rule{
		RuleNoPlaintextSecret,
		RuleNoLocalhostURL,
		RuleNoDefaultPassword,
	}
}

// Scan runs all rules against env and returns aggregated findings.
func Scan(env map[string]string, opts Options) []Finding {
	if len(opts.Rules) == 0 {
		opts.Rules = DefaultRules()
	}
	var findings []Finding
	for k, v := range env {
		for _, rule := range opts.Rules {
			findings = append(findings, rule(k, v)...)
		}
	}
	return findings
}

// RuleNoPlaintextSecret flags keys that look sensitive with short, common values.
func RuleNoPlaintextSecret(key, value string) []Finding {
	sensitivePatterns := []string{"SECRET", "PASSWORD", "TOKEN", "API_KEY", "PRIVATE"}
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) && len(value) > 0 && len(value) < 8 {
			return []Finding{{
				Key: key, Value: value,
				Severity: "warning",
				Message: fmt.Sprintf("key %q appears sensitive but has a suspiciously short value", key),
			}}
		}
	}
	return nil
}

// RuleNoLocalhostURL flags values containing localhost or 127.0.0.1.
func RuleNoLocalhostURL(key, value string) []Finding {
	if strings.Contains(value, "localhost") || strings.Contains(value, "127.0.0.1") {
		return []Finding{{
			Key: key, Value: value,
			Severity: "info",
			Message: fmt.Sprintf("key %q contains a localhost reference — may not be suitable for production", key),
		}}
	}
	return nil
}

// RuleNoDefaultPassword flags values that are common default passwords.
func RuleNoDefaultPassword(key, value string) []Finding {
	defaults := []string{"password", "secret", "changeme", "admin", "1234", "test"}
	lower := strings.ToLower(value)
	for _, d := range defaults {
		if lower == d {
			return []Finding{{
				Key: key, Value: value,
				Severity: "error",
				Message: fmt.Sprintf("key %q has a well-known default value %q", key, value),
			}}
		}
	}
	return nil
}
