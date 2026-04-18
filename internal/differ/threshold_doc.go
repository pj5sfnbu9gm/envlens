// Package differ — threshold.go
//
// ApplyThreshold filters a multi-target diff result set, retaining only those
// targets whose number of meaningful changes (added, removed, or changed keys)
// meets or exceeds a configurable minimum.
//
// This is useful when diffing many deployment targets and wanting to surface
// only the most divergent environments.
package differ
