// Package differ provides utilities for comparing environment variable maps.
//
// The Chain function computes sequential diffs across an ordered list of
// named environments (e.g. dev -> staging -> prod), returning a ChainResult
// for each consecutive pair. This is useful for auditing how configuration
// evolves across a promotion pipeline.
package differ
