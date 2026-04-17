// Package inspector provides key-level inspection of environment variables,
// returning metadata such as type guess, length, and sensitivity.
package inspector

import (
	"strconv"
	"strings"
)

// Entry holds inspection metadata for a single environment variable.
type Entry struct {
	Key       string
	Value     string
	Length    int
	Typeguess string
	Sensitive bool
	Empty     bool
}

// sensitivePatterns are substrings that suggest a key is sensitive.
var sensitivePatterns = []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS", "PRIVATE", "CREDENTIAL"}

// Inspect analyses each key/value pair in env and returns a slice of Entry.
func Inspect(env map[string]string) []Entry {
	entries := make([]Entry, 0, len(env))
	for k, v := range env {
		entries = append(entries, Entry{
			Key:       k,
			Value:     v,
			Length:    len(v),
			Typeguess: guessType(v),
			Sensitive: isSensitive(k),
			Empty:     v == "",
		})
	}
	return entries
}

func guessType(v string) string {
	if v == "" {
		return "empty"
	}
	if _, err := strconv.ParseBool(v); err == nil {
		return "bool"
	}
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return "int"
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return "float"
	}
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return "url"
	}
	return "string"
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
