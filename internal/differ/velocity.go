package differ

import (
	"sort"
)

// VelocityEntry records how rapidly a key's value is changing across a
// sequence of named environment snapshots.  The Rate field is the number
// of distinct transitions observed divided by the number of consecutive
// snapshot pairs examined.
type VelocityEntry struct {
	Key         string  `json:"key"`
	Transitions int     `json:"transitions"`
	Windows     int     `json:"windows"`
	Rate        float64 `json:"rate"` // transitions / windows
}

// VelocityOptions controls the behaviour of Velocity.
type VelocityOptions struct {
	// MinRate excludes entries whose rate is strictly below this value.
	// A value of 0 (the default) includes all keys that changed at least once.
	MinRate float64

	// TopN, when > 0, retains only the N entries with the highest rate.
	TopN int

	// IncludeStable, when true, also returns keys that never changed
	// (rate == 0) across the snapshot sequence.
	IncludeStable bool
}

// DefaultVelocityOptions returns a sensible default configuration.
func DefaultVelocityOptions() VelocityOptions {
	return VelocityOptions{
		MinRate:       0,
		TopN:          0,
		IncludeStable: false,
	}
}

// Velocity measures how frequently each key changes across an ordered
// sequence of named environment maps.  snapshots must be provided in
// chronological order; each entry is a (name, env) pair.  Consecutive
// pairs are diffed and transition counts are accumulated per key.
//
// The returned slice is sorted by Rate descending, then Key ascending.
func Velocity(snapshots []NamedEnv, opts VelocityOptions) []VelocityEntry {
	if len(snapshots) < 2 {
		return nil
	}

	transitions := map[string]int{}
	windowCount := len(snapshots) - 1

	for i := 0; i < len(snapshots)-1; i++ {
		a := snapshots[i].Env
		b := snapshots[i+1].Env

		// Collect all keys present in either snapshot.
		keys := unionKeys(a, b)
		for _, k := range keys {
			va, aOK := a[k]
			vb, bOK := b[k]
			if aOK != bOK || va != vb {
				transitions[k]++
			}
		}
	}

	// Build candidate set.
	allKeys := map[string]struct{}{}
	for _, s := range snapshots {
		for k := range s.Env {
			allKeys[k] = struct{}{}
		}
	}

	var entries []VelocityEntry
	for k := range allKeys {
		t := transitions[k]
		rate := float64(t) / float64(windowCount)
		if !opts.IncludeStable && t == 0 {
			continue
		}
		if rate < opts.MinRate {
			continue
		}
		entries = append(entries, VelocityEntry{
			Key:         k,
			Transitions: t,
			Windows:     windowCount,
			Rate:        rate,
		})
	}

	// Sort: highest rate first, then lexicographic key.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Rate != entries[j].Rate {
			return entries[i].Rate > entries[j].Rate
		}
		return entries[i].Key < entries[j].Key
	})

	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}

	return entries
}

// HasVelocityChanges returns true when at least one entry is present in
// the result, indicating at least one key changed across the snapshots.
func HasVelocityChanges(entries []VelocityEntry) bool {
	return len(entries) > 0
}

// unionKeys returns the deduplicated union of keys from two env maps.
func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
