package differ

import "sort"

// HeatmapEntry records how frequently a key has changed across a series of diffs.
type HeatmapEntry struct {
	Key       string
	Changes   int
	Targets   []string
}

// DefaultHeatmapOptions returns sensible defaults for Heatmap.
func DefaultHeatmapOptions() HeatmapOptions {
	return HeatmapOptions{
		MinChanges: 1,
		TopN:       0,
	}
}

// HeatmapOptions controls Heatmap behaviour.
type HeatmapOptions struct {
	// MinChanges filters out keys with fewer than this many changes.
	MinChanges int
	// TopN, if > 0, limits output to the N most-changed keys.
	TopN int
}

// Heatmap aggregates per-key change counts across multiple target result sets.
// results is a map of target-name → []Result (as produced by Diff or MultiDiff).
func Heatmap(results map[string][]Result, opts HeatmapOptions) []HeatmapEntry {
	type entry struct {
		changes int
		targets []string
	}
	agg := map[string]*entry{}

	for target, rs := range results {
		for _, r := range rs {
			if r.Status == StatusUnchanged {
				continue
			}
			e, ok := agg[r.Key]
			if !ok {
				e = &entry{}
				agg[r.Key] = e
			}
			e.changes++
			e.targets = append(e.targets, target)
		}
	}

	out := make([]HeatmapEntry, 0, len(agg))
	for key, e := range agg {
		if e.changes < opts.MinChanges {
			continue
		}
		sort.Strings(e.targets)
		out = append(out, HeatmapEntry{Key: key, Changes: e.changes, Targets: e.targets})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Changes != out[j].Changes {
			return out[i].Changes > out[j].Changes
		}
		return out[i].Key < out[j].Key
	})

	if opts.TopN > 0 && len(out) > opts.TopN {
		out = out[:opts.TopN]
	}
	return out
}

// HasHeatmapEntries returns true when at least one entry is present.
func HasHeatmapEntries(entries []HeatmapEntry) bool {
	return len(entries) > 0
}
