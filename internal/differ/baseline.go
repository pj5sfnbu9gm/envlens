package differ

// BaselineOptions controls how a baseline comparison is performed.
type BaselineOptions struct {
	// IgnoreUnchanged skips keys with identical values across all targets.
	IgnoreUnchanged bool
	// BaselineName is the name of the target treated as the reference.
	BaselineName string
}

// DefaultBaselineOptions returns sensible defaults.
func DefaultBaselineOptions() BaselineOptions {
	return BaselineOptions{
		IgnoreUnchanged: true,
		BaselineName:    "baseline",
	}
}

// CompareToBaseline diffs every non-baseline target against the named baseline
// target and returns a map of target name → diff results.
func CompareToBaseline(targets map[string]map[string]string, opts BaselineOptions) map[string][]Result {
	base, ok := targets[opts.BaselineName]
	if !ok {
		base = map[string]string{}
	}

	out := make(map[string][]Result)
	for name, env := range targets {
		if name == opts.BaselineName {
			continue
		}
		results := Diff(base, env)
		if opts.IgnoreUnchanged {
			filtered := results[:0]
			for _, r := range results {
				if r.Status != StatusUnchanged {
					filtered = append(filtered, r)
				}
			}
			results = filtered
		}
		out[name] = results
	}
	return out
}

// HasBaselineDifferences returns true if any target differs from the baseline.
func HasBaselineDifferences(results map[string][]Result) bool {
	for _, rs := range results {
		if len(rs) > 0 {
			return true
		}
	}
	return false
}
