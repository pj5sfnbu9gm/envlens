package differ

// CapOptions controls how Cap limits diff results per target.
type CapOptions struct {
	// MaxPerTarget is the maximum number of results to retain per target.
	// Results are ordered: Changed > Added > Removed > Unchanged.
	// A value of 0 means no limit.
	MaxPerTarget int

	// IgnoreUnchanged skips unchanged entries before applying the cap.
	IgnoreUnchanged bool
}

// DefaultCapOptions returns sensible defaults: no cap, ignore unchanged.
func DefaultCapOptions() CapOptions {
	return CapOptions{
		MaxPerTarget:    0,
		IgnoreUnchanged: true,
	}
}

// statusPriority assigns a sort priority to each DiffStatus for Cap ordering.
var statusPriority = map[DiffStatus]int{
	StatusChanged:   0,
	StatusAdded:     1,
	StatusRemoved:   2,
	StatusUnchanged: 3,
}

// Cap limits the number of diff results per target, prioritising the most
// significant changes (Changed > Added > Removed > Unchanged).
func Cap(results map[string][]Result, opts CapOptions) map[string][]Result {
	out := make(map[string][]Result, len(results))

	for target, entries := range results {
		filtered := entries
		if opts.IgnoreUnchanged {
			filtered = make([]Result, 0, len(entries))
			for _, r := range entries {
				if r.Status != StatusUnchanged {
					filtered = append(filtered, r)
				}
			}
		}

		// Stable priority sort: lower priority value = higher importance.
		sorted := make([]Result, len(filtered))
		copy(sorted, filtered)
		stableSortByPriority(sorted)

		if opts.MaxPerTarget > 0 && len(sorted) > opts.MaxPerTarget {
			sorted = sorted[:opts.MaxPerTarget]
		}

		out[target] = sorted
	}

	return out
}

// HasCapResults returns true when at least one target has results after capping.
func HasCapResults(results map[string][]Result) bool {
	for _, entries := range results {
		if len(entries) > 0 {
			return true
		}
	}
	return false
}

// stableSortByPriority performs an insertion sort (stable) by status priority.
func stableSortByPriority(rs []Result) {
	for i := 1; i < len(rs); i++ {
		for j := i; j > 0 && statusPriority[rs[j].Status] < statusPriority[rs[j-1].Status]; j-- {
			rs[j], rs[j-1] = rs[j-1], rs[j]
		}
	}
}
