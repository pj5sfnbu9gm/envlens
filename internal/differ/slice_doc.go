// Package differ provides utilities for comparing environment variable
// configurations across deployment targets.
//
// The Slice function partitions a flat list of diff Results into groups
// based on key prefixes (e.g. DB_, APP_, REDIS_). This is useful for
// presenting large diffs in a structured, readable way — one section
// per service or component boundary.
package differ
