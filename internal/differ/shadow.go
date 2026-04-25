package differ

// ShadowEntry represents a key that exists in a shadow (secondary) target
// but differs from or is absent in the primary target.
type ShadowEntry struct {
	Key          string
	PrimaryValue string
	ShadowValue  string
	OnlyInShadow bool
	OnlyInPrimary bool
}

// ShadowOptions controls Shadow behaviour.
type ShadowOptions struct {
	IncludeUnchanged bool
}

// DefaultShadowOptions returns sensible defaults.
func DefaultShadowOptions() ShadowOptions {
	return ShadowOptions{IncludeUnchanged: false}
}

// Shadow compares a primary environment against one or more shadow environments,
// returning per-key discrepancies. Keys present only in the primary are included
// with OnlyInPrimary set; keys present only in a shadow set OnlyInShadow.
func Shadow(primary map[string]string, shadows map[string]map[string]string, opts ShadowOptions) map[string][]ShadowEntry {
	result := make(map[string][]ShadowEntry)

	allKeys := make(map[string]struct{})
	for k := range primary {
		allKeys[k] = struct{}{}
	}
	for _, sh := range shadows {
		for k := range sh {
			allKeys[k] = struct{}{}
		}
	}

	for key := range allKeys {
		pv, inPrimary := primary[key]
		var entries []ShadowEntry
		for target, sh := range shadows {
			sv, inShadow := sh[key]
			if !inPrimary && !inShadow {
				continue
			}
			if !opts.IncludeUnchanged && inPrimary && inShadow && pv == sv {
				continue
			}
			_ = target
			entries = append(entries, ShadowEntry{
				Key:           key,
				PrimaryValue:  pv,
				ShadowValue:   sv,
				OnlyInShadow:  !inPrimary && inShadow,
				OnlyInPrimary: inPrimary && !inShadow,
			})
		}
		if len(entries) > 0 {
			result[key] = entries
		}
	}
	return result
}

// HasShadowDifferences returns true when Shadow produced any discrepancies.
func HasShadowDifferences(entries map[string][]ShadowEntry) bool {
	for _, v := range entries {
		if len(v) > 0 {
			return true
		}
	}
	return false
}
