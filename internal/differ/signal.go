package differ

import "sort"

// SignalEntry represents a key that appears as a significant signal
// across multiple diff results — i.e., it changed in many targets.
type SignalEntry struct {
	Key        string
	ChangeCount int
	Targets    []string
}

// SignalOptions controls how Signal filters and ranks entries.
type SignalOptions struct {
	// MinTargets is the minimum number of targets a key must change in.
	MinTargets int
	// TopN limits results to the top N keys by change count. 0 = no limit.
	TopN int
}

// DefaultSignalOptions returns sensible defaults.
func DefaultSignalOptions() SignalOptions {
	return SignalOptions{
		MinTargets: 2,
		TopN:       0,
	}
}

// Signal scans multi-target diff results and surfaces keys that changed
// across the most targets, acting as high-signal indicators of drift.
func Signal(targets map[string][]Result, opts SignalOptions) []SignalEntry {
	counts := map[string][]string{}

	for target, results := range targets {
		for _, r := range results {
			if r.Status == StatusChanged || r.Status == StatusAdded || r.Status == StatusRemoved {
				counts[r.Key] = append(counts[r.Key], target)
			}
		}
	}

	var entries []SignalEntry
	for key, tgts := range counts {
		if len(tgts) < opts.MinTargets {
			continue
		}
		sort.Strings(tgts)
		entries = append(entries, SignalEntry{
			Key:         key,
			ChangeCount: len(tgts),
			Targets:     tgts,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].ChangeCount != entries[j].ChangeCount {
			return entries[i].ChangeCount > entries[j].ChangeCount
		}
		return entries[i].Key < entries[j].Key
	})

	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}
	return entries
}

// HasSignals returns true if any signal entries were found.
func HasSignals(entries []SignalEntry) bool {
	return len(entries) > 0
}
