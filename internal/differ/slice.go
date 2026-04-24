package differ

// SliceResult holds the diff results for a single key slice (prefix group).
type SliceResult struct {
	Prefix  string
	Results []Result
}

// SliceOptions configures the Slice operation.
type SliceOptions struct {
	// Prefixes limits slicing to the given key prefixes.
	// If empty, all keys are grouped by their natural prefix.
	Prefixes []string
	// Delimiter separates prefix from the rest of the key.
	Delimiter string
	// IgnoreUnchanged skips unchanged results within each slice.
	IgnoreUnchanged bool
}

// DefaultSliceOptions returns sensible defaults for Slice.
func DefaultSliceOptions() SliceOptions {
	return SliceOptions{
		Delimiter:       "_",
		IgnoreUnchanged: true,
	}
}

// Slice partitions diff results into groups keyed by prefix.
// Each group contains only the results whose key starts with that prefix.
func Slice(results []Result, opts SliceOptions) []SliceResult {
	if len(results) == 0 {
		return nil
	}
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}

	order := []string{}
	buckets := map[string][]Result{}

	for _, r := range results {
		if opts.IgnoreUnchanged && r.Status == StatusUnchanged {
			continue
		}
		prefix := slicePrefix(r.Key, opts.Delimiter, opts.Prefixes)
		if _, seen := buckets[prefix]; !seen {
			order = append(order, prefix)
		}
		buckets[prefix] = append(buckets[prefix], r)
	}

	out := make([]SliceResult, 0, len(order))
	for _, p := range order {
		out = append(out, SliceResult{Prefix: p, Results: buckets[p]})
	}
	return out
}

// HasSliceChanges reports whether any SliceResult contains non-empty results.
func HasSliceChanges(slices []SliceResult) bool {
	for _, s := range slices {
		if len(s.Results) > 0 {
			return true
		}
	}
	return false
}

func slicePrefix(key, delimiter string, allowed []string) string {
	if len(allowed) > 0 {
		for _, p := range allowed {
			if len(key) > len(p) && key[:len(p)] == p {
				return p
			}
		}
		return "other"
	}
	for i := 0; i < len(key); i++ {
		if string(key[i]) == delimiter {
			return key[:i]
		}
	}
	return key
}
