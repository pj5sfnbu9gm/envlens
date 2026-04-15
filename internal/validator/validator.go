package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule applied to an environment map.
type Rule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// Finding represents a single validation issue.
type Finding struct {
	Key     string
	Rule    string
	Message string
}

var validKeyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// DefaultRules returns the built-in set of validation rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:    "required-value",
			Message: "key has an empty value",
			Check:   func(_, v string) bool { return strings.TrimSpace(v) == "" },
		},
		{
			Name:    "invalid-key-format",
			Message: "key does not match expected format (uppercase letters, digits, underscores)",
			Check:   func(k, _ string) bool { return !validKeyPattern.MatchString(k) },
		},
		{
			Name:    "leading-trailing-space",
			Message: "value has leading or trailing whitespace",
			Check:   func(_, v string) bool { return v != strings.TrimSpace(v) },
		},
	}
}

// Validate runs the provided rules against every key/value in env and
// returns a slice of findings (one per violated rule per key).
func Validate(env map[string]string, rules []Rule) []Finding {
	var findings []Finding
	for k, v := range env {
		for _, r := range rules {
			if r.Check(k, v) {
				findings = append(findings, Finding{
					Key:     k,
					Rule:    r.Name,
					Message: fmt.Sprintf("%s: %s", k, r.Message),
				})
			}
		}
	}
	return findings
}
