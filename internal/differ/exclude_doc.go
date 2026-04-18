// Package differ provides utilities for comparing environment variable maps
// across deployment targets.
//
// The Exclude function allows callers to strip specific keys or key prefixes
// from a set of diff results before reporting or further processing. This is
// useful for ignoring well-known volatile keys (e.g. timestamps, secrets) that
// are not meaningful to compare.
package differ
