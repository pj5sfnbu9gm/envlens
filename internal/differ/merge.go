package differ

// MergeResult holds the merged value for a key across targets.
type MergeResult struct {
	Key      string
	Value    string
	Sources  []string // target names that contributed this value
	Conflict bool     // true if targets disagreed
}

// MergeOptions controls how MergeDiff behaves.
type MergeOptions struct {
	PreferFirst bool // if true, use the first non-empty value; otherwise last wins
	SkipConflicts bool // if true, omit keys with conflicts from output
}

// DefaultMergeOptions returns sensible defaults.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		PreferFirst:   true,
		SkipConflicts: false,
	}
}

// MergeDiff collapses multiple target envs into a single merged view.
// Keys present in any target are included; conflicts are flagged.
func MergeDiff(targets map[string]map[string]string, opts MergeOptions) []MergeResult {
	if len(targets) == 0 {
		return nil
	}

	// collect all keys
	keySet := map[string]struct{}{}
	for _, env := range targets {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	names := make([]string, 0, len(targets))
	for name := range targets {
		names = append(names, name)
	}
	sortStrings(names)

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sortStrings(keys)

	var results []MergeResult
	for _, k := range keys {
		var chosen string
		var sources []string
		values := map[string]bool{}

		for _, name := range names {
			v, ok := targets[name][k]
			if !ok {
				continue
			}
			sources = append(sources, name)
			values[v] = true
			if opts.PreferFirst && chosen == "" {
				chosen = v
			} else if !opts.PreferFirst {
				chosen = v
			}
		}

		conflict := len(values) > 1
		if conflict && opts.SkipConflicts {
			continue
		}

		results = append(results, MergeResult{
			Key:      k,
			Value:    chosen,
			Sources:  sources,
			Conflict: conflict,
		})
	}
	return results
}

// HasMergeConflicts returns true if any result has a conflict.
func HasMergeConflicts(results []MergeResult) bool {
	for _, r := range results {
		if r.Conflict {
			return true
		}
	}
	return false
}
