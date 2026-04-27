package differ

// PruneOptions controls how Prune filters diff results.
type PruneOptions struct {
	// RemoveUnchanged drops all entries with status Unchanged.
	RemoveUnchanged bool
	// RemoveAdded drops all entries with status Added.
	RemoveAdded bool
	// RemoveRemoved drops all entries with status Removed.
	RemoveRemoved bool
	// RemoveChanged drops all entries with status Changed.
	RemoveChanged bool
	// Keys is an explicit set of keys to remove from results.
	Keys []string
	// Prefixes removes any entry whose key starts with one of these prefixes.
	Prefixes []string
}

// DefaultPruneOptions returns a PruneOptions that only removes unchanged entries.
func DefaultPruneOptions() PruneOptions {
	return PruneOptions{RemoveUnchanged: true}
}

// Prune filters a slice of Result values according to the given options.
// It returns a new slice containing only the entries that survive all filters.
func Prune(results []Result, opts PruneOptions) []Result {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make([]Result, 0, len(results))
	for _, r := range results {
		if opts.RemoveUnchanged && r.Status == StatusUnchanged {
			continue
		}
		if opts.RemoveAdded && r.Status == StatusAdded {
			continue
		}
		if opts.RemoveRemoved && r.Status == StatusRemoved {
			continue
		}
		if opts.RemoveChanged && r.Status == StatusChanged {
			continue
		}
		if _, found := keySet[r.Key]; found {
			continue
		}
		if pruneMatchesPrefix(r.Key, opts.Prefixes) {
			continue
		}
		out = append(out, r)
	}
	return out
}

// HasPruneResults reports whether any results remain after pruning.
func HasPruneResults(results []Result) bool {
	return len(results) > 0
}

func pruneMatchesPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
