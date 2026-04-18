package differ

// Projection filters a set of diff results to only include specified keys.
// This is useful when you want to focus on a subset of variables across targets.

// ProjectionOptions controls which keys are retained in the projection.
type ProjectionOptions struct {
	Keys       []string // exact keys to include
	Prefixes   []string // key prefixes to include
	Invert     bool     // if true, exclude matched keys instead
}

// DefaultProjectionOptions returns a ProjectionOptions with no filters applied.
func DefaultProjectionOptions() ProjectionOptions {
	return ProjectionOptions{}
}

// Project filters DiffResults to only include entries matching the given options.
func Project(results []DiffResult, opts ProjectionOptions) []DiffResult {
	if len(opts.Keys) == 0 && len(opts.Prefixes) == 0 {
		return results
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	var out []DiffResult
	for _, r := range results {
		matched := keySet[r.Key] || matchesPrefix(r.Key, opts.Prefixes)
		if opts.Invert {
			matched = !matched
		}
		if matched {
			out = append(out, r)
		}
	}
	return out
}

// HasProjectedChanges returns true if any projected result is not Unchanged.
func HasProjectedChanges(results []DiffResult) bool {
	for _, r := range results {
		if r.Status != StatusUnchanged {
			return true
		}
	}
	return false
}

func matchesPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
