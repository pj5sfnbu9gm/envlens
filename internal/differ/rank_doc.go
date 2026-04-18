// Package differ — rank.go
//
// Rank aggregates a MultiDiff result set and counts how many deployment
// targets each key changed in. This makes it easy to identify the most
// volatile keys across an environment fleet.
//
// Usage:
//
//	opts := differ.DefaultRankOptions()
//	opts.TopN = 10
//	entries := differ.Rank(multiResults, opts)
//	for _, e := range entries {
//		fmt.Printf("%s changed in %d targets\n", e.Key, e.Changes)
//	}
package differ
