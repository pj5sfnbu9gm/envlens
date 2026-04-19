package differ

import "sort"

// RollupEntry summarizes change counts across all targets for a single key.
type RollupEntry struct {
	Key       string
	Added     int
	Removed   int
	Changed   int
	Unchanged int
	Total     int
}

// RollupOptions controls Rollup behaviour.
type RollupOptions struct {
	// MinChanges filters out keys whose total change count is below this value.
	MinChanges int
}

// DefaultRollupOptions returns sensible defaults.
func DefaultRollupOptions() RollupOptions {
	return RollupOptions{MinChanges: 1}
}

// Rollup aggregates per-target diff results into a per-key summary.
// targets is a map of target name -> slice of Result for that target.
func Rollup(targets map[string][]Result, opts RollupOptions) []RollupEntry {
	counts := map[string]*RollupEntry{}

	for _, results := range targets {
		for _, r := range results {
			e, ok := counts[r.Key]
			if !ok {
				e = &RollupEntry{Key: r.Key}
				counts[r.Key] = e
			}
			switch r.Status {
			case StatusAdded:
				e.Added++
			case StatusRemoved:
				e.Removed++
			case StatusChanged:
				e.Changed++
			case StatusUnchanged:
				e.Unchanged++
			}
			e.Total = e.Added + e.Removed + e.Changed
		}
	}

	var out []RollupEntry
	for _, e := range counts {
		if e.Total >= opts.MinChanges {
			out = append(out, *e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Total != out[j].Total {
			return out[i].Total > out[j].Total
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// HasRollupChanges returns true if any entry has at least one change.
func HasRollupChanges(entries []RollupEntry) bool {
	for _, e := range entries {
		if e.Total > 0 {
			return true
		}
	}
	return false
}
