// Package differ provides utilities for comparing environment variable
// configurations across deployment targets.
//
// Heatmap aggregates per-key change frequency across multiple diff result
// sets, making it easy to identify which variables change most often.
// Results are sorted by change count descending so the "hottest" keys
// appear first. Use HeatmapOptions.TopN to cap the output and
// HeatmapOptions.MinChanges to filter out infrequently changed keys.
package differ
