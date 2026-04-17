// Package stringer provides utilities for converting environment maps
// into various string representations.
package stringer

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how the environment map is stringified.
type Options struct {
	Sorted    bool
	Separator string // between key and value, default "="
	Delimiter string // between entries, default "\n"
	QuoteValues bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Sorted:    true,
		Separator: "=",
		Delimiter: "\n",
		QuoteValues: false,
	}
}

// Stringify converts an env map to a single string.
func Stringify(env map[string]string, opts Options) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.Sorted {
		sort.Strings(keys)
	}

	sep := opts.Separator
	if sep == "" {
		sep = "="
	}
	delim := opts.Delimiter
	if delim == "" {
		delim = "\n"
	}

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		if opts.QuoteValues {
			v = fmt.Sprintf("%q", v)
		}
		parts = append(parts, k+sep+v)
	}
	return strings.Join(parts, delim)
}

// ToLines returns each key=value pair as a slice of strings.
func ToLines(env map[string]string, opts Options) []string {
	result := Stringify(env, opts)
	if result == "" {
		return []string{}
	}
	return strings.Split(result, opts.Delimiter)
}
