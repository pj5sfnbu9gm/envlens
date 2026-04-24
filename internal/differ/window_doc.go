// Package differ — window.go
//
// Window applies a sliding window over an ordered sequence of named environment
// snapshots, producing a diff for each consecutive group of Size snapshots.
//
// This is useful for detecting gradual configuration drift across a series of
// deployments or time-stamped captures, where you want to see how the
// environment changed between adjacent releases rather than comparing only the
// first and last.
package differ
