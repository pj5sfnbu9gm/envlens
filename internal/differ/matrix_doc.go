// Package differ provides utilities for comparing environment variable maps.
//
// The Matrix function computes pairwise diffs across all provided named targets,
// returning a flat list of MatrixEntry values — one per ordered pair (A→B).
// This is useful for visualising how every environment relates to every other.
package differ
