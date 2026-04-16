// Package grouper groups environment variables by a common prefix delimiter.
package grouper

import "strings"

// Options controls grouper behaviour.
type Options struct {
	// Delimiter separates the prefix from the rest of the key (default "_").
	Delimiter string
	// MinGroupSize skips groups with fewer members than this value (0 = include all).
	MinGroupSize int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Delimiter: "_", MinGroupSize: 0}
}

// Group partitions env into named groups keyed by their first prefix segment.
// Keys with no delimiter are placed in the "" (empty-string) group.
func Group(env map[string]string, opts Options) map[string]map[string]string {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}

	raw := make(map[string]map[string]string)
	for k, v := range env {
		prefix := ""
		if idx := strings.Index(k, opts.Delimiter); idx > 0 {
			prefix = k[:idx]
		}
		if raw[prefix] == nil {
			raw[prefix] = make(map[string]string)
		}
		raw[prefix][k] = v
	}

	if opts.MinGroupSize <= 1 {
		return raw
	}

	filtered := make(map[string]map[string]string)
	for prefix, members := range raw {
		if len(members) >= opts.MinGroupSize {
			filtered[prefix] = members
		}
	}
	return filtered
}
