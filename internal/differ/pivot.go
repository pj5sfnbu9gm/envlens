package differ

// PivotEntry represents a single key's values across multiple targets.
type PivotEntry struct {
	Key     string
	Values  map[string]string // target name -> value
	Uniform bool             // true if all targets agree
}

// PivotOptions controls PivotDiff behaviour.
type PivotOptions struct {
	IncludeUnchanged bool
}

// DefaultPivotOptions returns sensible defaults.
func DefaultPivotOptions() PivotOptions {
	return PivotOptions{IncludeUnchanged: false}
}

// PivotDiff reorganises per-target diff results into a key-centric view.
// targets is a map of target-name -> env map.
func PivotDiff(targets map[string]map[string]string, opts PivotOptions) []PivotEntry {
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
	for n := range targets {
		names = append(names, n)
	}
	sortStrings(names)

	allKeys := make([]string, 0, len(keySet))
	for k := range keySet {
		allKeys = append(allKeys, k)
	}
	sortStrings(allKeys)

	var entries []PivotEntry
	for _, key := range allKeys {
		values := make(map[string]string, len(targets))
		for _, name := range names {
			if v, ok := targets[name][key]; ok {
				values[name] = v
			}
		}

		uniform := true
		var first string
		set := false
		for _, name := range names {
			v := values[name]
			if !set {
				first = v
				set = true
			} else if v != first {
				uniform = false
				break
			}
		}

		if !opts.IncludeUnchanged && uniform {
			continue
		}

		entries = append(entries, PivotEntry{
			Key:     key,
			Values:  values,
			Uniform: uniform,
		})
	}
	return entries
}

// HasPivotDifferences returns true if any entry is non-uniform.
func HasPivotDifferences(entries []PivotEntry) bool {
	for _, e := range entries {
		if !e.Uniform {
			return true
		}
	}
	return false
}
