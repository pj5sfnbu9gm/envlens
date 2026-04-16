// Package redactor replaces sensitive environment variable values with a
// configurable placeholder string.
//
// Keys are matched for redaction either by exact name or by prefix. The
// original map is never mutated; a new map is returned inside Result.
//
// Example:
//
//	opts := redactor.DefaultOptions()
//	opts.Keys = append(opts.Keys, "DB_PASSWORD")
//	result := redactor.Redact(env, opts)
//	fmt.Println(result.Redacted["DB_PASSWORD"]) // [REDACTED]
package redactor
