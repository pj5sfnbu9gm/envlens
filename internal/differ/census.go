package differ

import "sort"

// CensusEntry holds the count of how many targets define a given key.
type CensusEntry struct {
	Key      string
	Count    int
	Targets  []string
	Coverage float64 // fraction of targets that define this key
}

// CensusOptions controls Census behaviour.
type CensusOptions struct {
	// MinCoverage filters out entries whose coverage is below this value (0–1).
	MinCoverage float64
	// ExcludeUniversal drops keys that appear in every target.
	ExcludeUniversal bool
}

// DefaultCensusOptions returns sensible defaults.
func DefaultCensusOptions() CensusOptions {
	return CensusOptions{
		MinCoverage:      0,
		ExcludeUniversal: false,
	}
}

// Census counts, across all targets, how many define each key.
// targets is a map of target-name → env map.
func Census(targets map[string]map[string]string, opts CensusOptions) []CensusEntry {
	if len(targets) == 0 {
		return nil
	}

	total := float64(len(targets))
	counts := map[string][]string{}

	for name, env := range targets {
		for k := range env {
			counts[k] = append(counts[k], name)
		}
	}

	var out []CensusEntry
	for k, tgts := range counts {
		cov := float64(len(tgts)) / total
		if cov < opts.MinCoverage {
			continue
		}
		if opts.ExcludeUniversal && len(tgts) == len(targets) {
			continue
		}
		sort.Strings(tgts)
		out = append(out, CensusEntry{
			Key:      k,
			Count:    len(tgts),
			Targets:  tgts,
			Coverage: cov,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// HasCensusGaps reports whether any key is missing from at least one target.
func HasCensusGaps(entries []CensusEntry, totalTargets int) bool {
	for _, e := range entries {
		if e.Count < totalTargets {
			return true
		}
	}
	return false
}
