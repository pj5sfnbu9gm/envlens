package differ

// ThresholdOptions controls how the threshold filter behaves.
type ThresholdOptions struct {
	// MinChanges is the minimum number of changed/added/removed keys required
	// for a target to be included in the output.
	MinChanges int
	// IncludeUnchanged includes targets that meet the threshold but have
	// unchanged keys as well.
	IncludeUnchanged bool
}

// DefaultThresholdOptions returns sensible defaults.
func DefaultThresholdOptions() ThresholdOptions {
	return ThresholdOptions{
		MinChanges:       1,
		IncludeUnchanged: false,
	}
}

// ThresholdResult holds filtered diff results for a single target.
type ThresholdResult struct {
	Target  string
	Results []Result
	Count   int
}

// ApplyThreshold filters a map of target→results, keeping only targets whose
// number of non-unchanged entries meets or exceeds opts.MinChanges.
func ApplyThreshold(targets map[string][]Result, opts ThresholdOptions) []ThresholdResult {
	var out []ThresholdResult
	keys := sortStrings(mapKeys(targets))
	for _, target := range keys {
		results := targets[target]
		count := 0
		for _, r := range results {
			if r.Status != StatusUnchanged {
				count++
			}
		}
		if count < opts.MinChanges {
			continue
		}
		filtered := results
		if !opts.IncludeUnchanged {
			filtered = make([]Result, 0, count)
			for _, r := range results {
				if r.Status != StatusUnchanged {
					filtered = append(filtered, r)
				}
			}
		}
		out = append(out, ThresholdResult{Target: target, Results: filtered, Count: count})
	}
	return out
}

// HasThresholdResults returns true if any target survived the threshold filter.
func HasThresholdResults(results []ThresholdResult) bool {
	return len(results) > 0
}

func mapKeys(m map[string][]Result) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
