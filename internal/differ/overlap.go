package differ

// OverlapResult holds keys that exist in some but not all targets.
type OverlapResult struct {
	Key            string
	PresentIn      []string
	AbsentFrom     []string
}

// FindOverlap returns keys that appear in at least one but not all targets.
// Keys present in every target are excluded (use Intersect for those).
func FindOverlap(targets map[string]map[string]string) []OverlapResult {
	if len(targets) == 0 {
		return nil
	}

	// Collect all keys across targets.
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

	var results []OverlapResult
	for _, k := range keys {
		var present, absent []string
		for _, name := range names {
			if _, ok := targets[name][k]; ok {
				present = append(present, name)
			} else {
				absent = append(absent, name)
			}
		}
		// Only include keys missing from at least one target.
		if len(absent) > 0 {
			results = append(results, OverlapResult{
				Key:        k,
				PresentIn:  present,
				AbsentFrom: absent,
			})
		}
	}
	return results
}

// HasOverlap returns true if any key is missing from at least one target.
func HasOverlap(targets map[string]map[string]string) bool {
	return len(FindOverlap(targets)) > 0
}
