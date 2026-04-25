package differ

import (
	"sort"
)

// BloomEntry represents a key's presence across a set of named targets.
type BloomEntry struct {
	Key      string
	PresentIn []string
	AbsentIn  []string
}

// BloomOptions controls Bloom behaviour.
type BloomOptions struct {
	// MinPresence filters entries that appear in fewer than MinPresence targets.
	MinPresence int
}

// DefaultBloomOptions returns sensible defaults.
func DefaultBloomOptions() BloomOptions {
	return BloomOptions{MinPresence: 1}
}

// Bloom analyses which keys are present or absent across all targets.
// It returns one BloomEntry per unique key found in any target.
func Bloom(targets map[string]map[string]string, opts BloomOptions) []BloomEntry {
	if len(targets) == 0 {
		return nil
	}

	allKeys := map[string]struct{}{}
	for _, env := range targets {
		for k := range env {
			allKeys[k] = struct{}{}
		}
	}

	names := make([]string, 0, len(targets))
	for name := range targets {
		names = append(names, name)
	}
	sort.Strings(names)

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var results []BloomEntry
	for _, k := range keys {
		var present, absent []string
		for _, name := range names {
			if _, ok := targets[name][k]; ok {
				present = append(present, name)
			} else {
				absent = append(absent, name)
			}
		}
		if len(present) < opts.MinPresence {
			continue
		}
		results = append(results, BloomEntry{
			Key:       k,
			PresentIn: present,
			AbsentIn:  absent,
		})
	}
	return results
}

// HasBloomGaps returns true if any key is absent from at least one target.
func HasBloomGaps(entries []BloomEntry) bool {
	for _, e := range entries {
		if len(e.AbsentIn) > 0 {
			return true
		}
	}
	return false
}
