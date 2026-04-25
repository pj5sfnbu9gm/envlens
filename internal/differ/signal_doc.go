// Package differ — signal.go
//
// Signal surfaces high-signal environment variable keys that changed
// across the most deployment targets in a multi-target diff run.
//
// Use Signal when you want to quickly identify which keys are drifting
// broadly across your fleet rather than in a single isolated target.
//
// Example:
//
//	entries := differ.Signal(multiResults, differ.DefaultSignalOptions())
//	for _, e := range entries {
//		fmt.Printf("%s changed in %d targets\n", e.Key, e.ChangeCount)
//	}
package differ
