// Package filter provides utilities for filtering environment variable maps
// based on key patterns, prefixes, and custom predicates.
package filter

import (
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	// Prefixes limits results to keys starting with any of these prefixes.
	Prefixes []string
	// Contains limits results to keys containing this substring (case-insensitive).
	Contains string
	// ExcludeKeys removes specific keys from the result.
	ExcludeKeys []string
}

// Apply filters the given env map according to the provided Options.
// It returns a new map containing only the entries that pass all active filters.
func Apply(env map[string]string, opts Options) map[string]string {
	result := make(map[string]string, len(env))

	excluded := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = struct{}{}
	}

	for k, v := range env {
		if _, skip := excluded[k]; skip {
			continue
		}
		if len(opts.Prefixes) > 0 && !matchesAnyPrefix(k, opts.Prefixes) {
			continue
		}
		if opts.Contains != "" && !strings.Contains(strings.ToUpper(k), strings.ToUpper(opts.Contains)) {
			continue
		}
		result[k] = v
	}

	return result
}

// matchesAnyPrefix returns true if key starts with at least one of the given prefixes.
func matchesAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
