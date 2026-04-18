package differ

// Squash merges a sequence of diff result sets into a single consolidated
// result set, keeping the most significant change status for each key.
// Priority order: Changed > Added > Removed > Unchanged.

// SquashResult holds the consolidated result for a key across multiple diffs.
type SquashResult struct {
	Key    string
	Status string // "added", "removed", "changed", "unchanged"
	Old    string
	New    string
}

// Squash consolidates multiple []Result slices into a single []SquashResult.
func Squash(sets ...[]Result) []SquashResult {
	priority := map[string]int{
		"unchanged": 0,
		"removed":   1,
		"added":     2,
		"changed":   3,
	}

	best := make(map[string]SquashResult)

	for _, results := range sets {
		for _, r := range results {
			existing, ok := best[r.Key]
			if !ok || priority[r.Status] > priority[existing.Status] {
				best[r.Key] = SquashResult{
					Key:    r.Key,
					Status: r.Status,
					Old:    r.Old,
					New:    r.New,
				}
			}
		}
	}

	keys := make([]string, 0, len(best))
	for k := range best {
		keys = append(keys, k)
	}
	sortStrings(keys)

	out := make([]SquashResult, 0, len(keys))
	for _, k := range keys {
		out = append(out, best[k])
	}
	return out
}

// HasSquashedChanges returns true if any squashed result is not "unchanged".
func HasSquashedChanges(results []SquashResult) bool {
	for _, r := range results {
		if r.Status != "unchanged" {
			return true
		}
	}
	return false
}
