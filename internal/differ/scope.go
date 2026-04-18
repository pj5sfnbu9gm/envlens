package differ

// ScopeOptions controls which keys are included in a scoped diff.
type ScopeOptions struct {
	// Keys is an explicit list of keys to include.
	Keys []string
	// Prefixes limits the diff to keys matching any of the given prefixes.
	Prefixes []string
	// IgnoreUnchanged skips unchanged results from the output.
	IgnoreUnchanged bool
}

// ScopeResult holds the diff result for a single key within a scope.
type ScopeResult struct {
	Key    string
	Status string
	Old    string
	New    string
}

// Scope runs a diff between base and target, restricted to the keys and
// prefixes defined in opts.
func Scope(base, target map[string]string, opts ScopeOptions) []ScopeResult {
	all := Diff(base, target)
	var out []ScopeResult
	for _, r := range all {
		if !scopeMatches(r.Key, opts) {
			continue
		}
		if opts.IgnoreUnchanged && r.Status == "unchanged" {
			continue
		}
		out = append(out, ScopeResult{
			Key:    r.Key,
			Status: r.Status,
			Old:    r.Old,
			New:    r.New,
		})
	}
	return out
}

// HasScopeChanges returns true if any scoped result is not unchanged.
func HasScopeChanges(results []ScopeResult) bool {
	for _, r := range results {
		if r.Status != "unchanged" {
			return true
		}
	}
	return false
}

func scopeMatches(key string, opts ScopeOptions) bool {
	if len(opts.Keys) == 0 && len(opts.Prefixes) == 0 {
		return true
	}
	for _, k := range opts.Keys {
		if k == key {
			return true
		}
	}
	for _, p := range opts.Prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
