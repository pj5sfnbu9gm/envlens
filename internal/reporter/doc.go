// Package reporter provides formatting and output capabilities for envlens
// diff results. It supports multiple output formats including human-readable
// text (with optional ANSI color) and structured JSON.
//
// Basic usage:
//
//	opts := reporter.DefaultOptions()
//	err := reporter.Report(results, "staging", "production", opts)
//
// To output JSON:
//
//	opts := reporter.Options{
//		Format: reporter.FormatJSON,
//		Writer: os.Stdout,
//	}
//	err := reporter.Report(results, "staging", "production", opts)
package reporter
