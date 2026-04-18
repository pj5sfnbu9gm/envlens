package differ

import "sort"

// RankEntry holds a key and the number of targets in which it changed.
type RankEntry struct {
	Key     string
	Changes int
}

// RankOptions controls Rank behaviour.
type RankOptions struct {
	// TopN limits results to the N most-changed keys. 0 means no limit.
	TopN int
	// MinChanges filters out keys with fewer changes than this threshold.
	MinChanges int
}

// DefaultRankOptions returns sensible defaults.
func DefaultRankOptions() RankOptions {
	return RankOptions{TopN: 0, MinChanges: 1}
}

// Rank counts how many targets each key changed across a MultiDiff result and
// returns entries sorted by change count descending, then key ascending.
func Rank(results map[string][]Result, opts RankOptions) []RankEntry {
	counts := map[string]int{}
	for _, res := range results {
		for _, r := range res {
			if r.Status != StatusUnchanged {
				counts[r.Key]++
			}
		}
	}

	var entries []RankEntry
	for k, c := range counts {
		if c >= opts.MinChanges {
			entries = append(entries, RankEntry{Key: k, Changes: c})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Changes != entries[j].Changes {
			return entries[i].Changes > entries[j].Changes
		}
		return entries[i].Key < entries[j].Key
	})

	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}
	return entries
}

// HasRankResults returns true when at least one entry is present.
func HasRankResults(entries []RankEntry) bool {
	return len(entries) > 0
}
