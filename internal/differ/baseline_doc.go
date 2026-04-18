// Package differ provides utilities for comparing environment variable maps.
//
// The baseline sub-feature (baseline.go) allows comparing multiple deployment
// targets against a single named reference target (the "baseline"). This is
// useful for auditing drift between environments such as staging and production
// relative to a known-good local or CI configuration.
//
// Example:
//
//	results := differ.CompareToBaseline(targets, differ.DefaultBaselineOptions())
//	if differ.HasBaselineDifferences(results) {
//		// report or fail
//	}
package differ
