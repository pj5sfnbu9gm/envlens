// Package flattener converts nested key structures (e.g. APP__DB__HOST) into
// dot-notation or other flat representations, and vice versa.
package flattener

import "strings"

// Options controls flattening behaviour.
type Options struct {
	// Separator is the string used to detect / produce nesting levels.
	// Defaults to "__".
	Separator string
	// OutputSeparator replaces Separator in the output key.
	// Defaults to ".".
	OutputSeparator string
	// SkipEmpty drops entries whose value is the empty string.
	SkipEmpty bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator:       "__",
		OutputSeparator: ".",
	}
}

// FlattenResult holds a single flattened entry.
type FlattenResult struct {
	OriginalKey string
	FlatKey     string
	Value       string
}

// Flatten converts keys that contain Separator into OutputSeparator-delimited
// flat keys. Keys without the separator are passed through unchanged.
func Flatten(env map[string]string, opts Options) []FlattenResult {
	if opts.Separator == "" {
		opts.Separator = "__"
	}
	if opts.OutputSeparator == "" {
		opts.OutputSeparator = "."
	}

	results := make([]FlattenResult, 0, len(env))
	for k, v := range env {
		if opts.SkipEmpty && v == "" {
			continue
		}
		flat := strings.ReplaceAll(k, opts.Separator, opts.OutputSeparator)
		results = append(results, FlattenResult{
			OriginalKey: k,
			FlatKey:     flat,
			Value:       v,
		})
	}
	return results
}

// ToMap converts a slice of FlattenResult into a map keyed by FlatKey.
func ToMap(results []FlattenResult) map[string]string {
	out := make(map[string]string, len(results))
	for _, r := range results {
		out[r.FlatKey] = r.Value
	}
	return out
}
