// Package normalizer provides utilities for normalizing environment variable
// keys and values according to configurable rules.
//
// Supported normalizations include:
//   - Uppercasing keys
//   - Trimming leading/trailing whitespace from keys and values
//   - Replacing hyphens with underscores in keys
//   - Removing entries with empty values
//
// Each normalization operation returns both the transformed env map and a
// slice of Result records describing what changed, enabling audit trails
// and reporting downstream.
package normalizer
