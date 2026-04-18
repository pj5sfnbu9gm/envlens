package differ

// UniqueResult holds keys that are exclusive to a single target.
type UniqueResult struct {
	Target string
	Keys   []string
}

// FindUnique returns, for each target, the set of keys that exist only in
// that target and in none of the others.
func FindUnique(targets map[string]map[string]string) []UniqueResult {
	if len(targets) == 0 {
		return nil
	}

	// Count how many targets each key appears in.
	occurrences := make(map[string]int)
	for _, env := range targets {
		for k := range env {
			occurrences[k]++
		}
	}

	names := make([]string, 0, len(targets))
	for name := range targets {
		names = append(names, name)
	}
	sortStrings(names)

	results := make([]UniqueResult, 0, len(names))
	for _, name := range names {
		env := targets[name]
		var unique []string
		for k := range env {
			if occurrences[k] == 1 {
				unique = append(unique, k)
			}
		}
		sortStrings(unique)
		results = append(results, UniqueResult{Target: name, Keys: unique})
	}
	return results
}

// HasUniqueKeys returns true if any target contains keys not present in others.
func HasUniqueKeys(results []UniqueResult) bool {
	for _, r := range results {
		if len(r.Keys) > 0 {
			return true
		}
	}
	return false
}
