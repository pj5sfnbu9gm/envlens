package differ

// ExcludeOptions controls which keys are excluded from diff results.
type ExcludeOptions struct {
	// Keys is a list of exact key names to exclude.
	Keys []string
	// Prefixes is a list of key prefixes to exclude.
	Prefixes []string
}

// Exclude filters out diff results whose keys match the given options.
// It returns a new slice with matching entries removed.
func Exclude(results []Result, opts ExcludeOptions) []Result {
	if len(opts.Keys) == 0 && len(opts.Prefixes) == 0 {
		return results
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	filtered := make([]Result, 0, len(results))
	for _, r := range results {
		if _, ok := keySet[r.Key]; ok {
			continue
		}
		if hasPrefix(r.Key, opts.Prefixes) {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

// HasExcluded returns true if any results would be removed by the given options.
func HasExcluded(results []Result, opts ExcludeOptions) bool {
	return len(Exclude(results, opts)) < len(results)
}

func hasPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
