// Package trimmer provides utilities for cleaning environment variable maps
// by removing keys or values that match configurable criteria such as empty
// values, blank-only values, or keys matching specific patterns.
package trimmer

import "strings"

// Options controls which entries are removed during trimming.
type Options struct {
	// RemoveEmpty removes keys whose values are the empty string.
	RemoveEmpty bool
	// RemoveBlank removes keys whose values contain only whitespace.
	RemoveBlank bool
	// RemovePrefixes removes keys that start with any of the given prefixes.
	RemovePrefixes []string
	// RemoveKeys removes exact key matches from this set.
	RemoveKeys []string
}

// DefaultOptions returns an Options with RemoveEmpty and RemoveBlank enabled.
func DefaultOptions() Options {
	return Options{
		RemoveEmpty: true,
		RemoveBlank: true,
	}
}

// Trim returns a new map with entries removed according to opts.
// The original map is not modified.
func Trim(env map[string]string, opts Options) map[string]string {
	exactRemove := make(map[string]struct{}, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		exactRemove[k] = struct{}{}
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if opts.RemoveEmpty && v == "" {
			continue
		}
		if opts.RemoveBlank && strings.TrimSpace(v) == "" {
			continue
		}
		if _, ok := exactRemove[k]; ok {
			continue
		}
		if hasAnyPrefix(k, opts.RemovePrefixes) {
			continue
		}
		result[k] = v
	}
	return result
}

func hasAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
