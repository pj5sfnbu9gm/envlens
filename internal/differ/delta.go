package differ

import "sort"

// DeltaEntry represents the net change summary for a single key across a
// sequence of named environments (e.g. dev → staging → prod).
type DeltaEntry struct {
	Key      string
	Adds     int
	Removals int
	Changes  int
	Net      int // Adds - Removals
}

// DeltaOptions controls Delta behaviour.
type DeltaOptions struct {
	// MinNet filters entries whose absolute net value is below this threshold.
	MinNet int
	// IncludeZero includes entries with no net change when true.
	IncludeZero bool
}

// DefaultDeltaOptions returns sensible defaults.
func DefaultDeltaOptions() DeltaOptions {
	return DeltaOptions{MinNet: 0, IncludeZero: false}
}

// Delta aggregates per-key add/remove/change counts across multiple
// MultiDiff result maps and returns a sorted slice of DeltaEntry.
func Delta(targets map[string][]Result, opts DeltaOptions) []DeltaEntry {
	counts := map[string]*DeltaEntry{}

	for _, results := range targets {
		for _, r := range results {
			e, ok := counts[r.Key]
			if !ok {
				e = &DeltaEntry{Key: r.Key}
				counts[r.Key] = e
			}
			switch r.Status {
			case StatusAdded:
				e.Adds++
			case StatusRemoved:
				e.Removals++
			case StatusChanged:
				e.Changes++
			}
		}
	}

	var out []DeltaEntry
	for _, e := range counts {
		e.Net = e.Adds - e.Removals
		net := e.Net
		if net < 0 {
			net = -net
		}
		if net < opts.MinNet {
			continue
		}
		if !opts.IncludeZero && e.Adds == 0 && e.Removals == 0 && e.Changes == 0 {
			continue
		}
		out = append(out, *e)
	}

	sort.Slice(out, func(i, j int) bool {
		ni, nj := out[i].Net, out[j].Net
		if ni < 0 {
			ni = -ni
		}
		if nj < 0 {
			nj = -nj
		}
		if ni != nj {
			return ni > nj
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// HasDeltaEntries returns true if any entries exist.
func HasDeltaEntries(entries []DeltaEntry) bool { return len(entries) > 0 }
