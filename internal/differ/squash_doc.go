// Package differ provides utilities for comparing environment variable maps.
//
// Squash consolidates multiple diff result sets into a single authoritative
// result set by selecting the highest-priority status for each key.
// This is useful when diffing across a chain of environments and you want
// a single view of what changed end-to-end.
package differ
