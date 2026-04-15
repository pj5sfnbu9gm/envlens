package auditor

import (
	"fmt"
	"sort"
)

// Rule defines a validation rule applied to environment variable keys or values.
type Rule struct {
	Name    string
	Check   func(key, value string) error
}

// Finding represents a single audit finding for a key.
type Finding struct {
	Key     string
	Rule    string
	Message string
}

// Audit runs all provided rules against the given env map and returns findings.
func Audit(env map[string]string, rules []Rule) []Finding {
	var findings []Finding

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := env[key]
		for _, rule := range rules {
			if err := rule.Check(key, val); err != nil {
				findings = append(findings, Finding{
					Key:     key,
					Rule:    rule.Name,
					Message: err.Error(),
				})
			}
		}
	}

	return findings
}

// DefaultRules returns a set of built-in audit rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name: "no-empty-value",
			Check: func(key, value string) error {
				if value == "" {
					return fmt.Errorf("key %q has an empty value", key)
				}
				return nil
			},
		},
		{
			Name: "no-whitespace-key",
			Check: func(key, value string) error {
				for _, ch := range key {
					if ch == ' ' || ch == '\t' {
						return fmt.Errorf("key %q contains whitespace", key)
					}
				}
				return nil
			},
		},
		{
			Name: "uppercase-key",
			Check: func(key, value string) error {
				for _, ch := range key {
					if ch >= 'a' && ch <= 'z' {
						return fmt.Errorf("key %q contains lowercase letters", key)
					}
				}
				return nil
			},
		},
	}
}
