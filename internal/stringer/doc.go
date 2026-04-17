// Package stringer converts environment variable maps into human-readable
// or machine-parseable string formats.
//
// It supports configurable key-value separators, entry delimiters,
// optional value quoting, and sorted output for deterministic results.
//
// Example:
//
//	env := map[string]string{"APP": "envlens", "PORT": "8080"}
//	opts := stringer.DefaultOptions()
//	output := stringer.Stringify(env, opts)
//	// APP=envlens
//	// PORT=8080
package stringer
