package differ

import "sort"

// GradientEntry represents the change intensity for a single key across
// an ordered sequence of targets, tracking how a value evolves over time.
type GradientEntry struct {
	Key       string
	Steps     []string // target names in order
	Values    []string // value at each step (empty string if absent)
	Changes   int      // number of step-to-step transitions that differ
	Direction string   // "up", "down", "stable", "volatile"
}

// GradientOptions controls Gradient behaviour.
type GradientOptions struct {
	MinChanges     int
	IncludeStable  bool
}

// DefaultGradientOptions returns sensible defaults.
func DefaultGradientOptions() GradientOptions {
	return GradientOptions{
		MinChanges:    1,
		IncludeStable: false,
	}
}

// Gradient analyses how each key's value shifts across an ordered slice of
// named targets. Unlike Trend (which works on windows of diffs), Gradient
// operates directly on raw env maps so callers can supply any ordered set.
func Gradient(targets []NamedEnv, opts GradientOptions) []GradientEntry {
	if len(targets) == 0 {
		return nil
	}

	// collect all keys
	keySet := map[string]struct{}{}
	for _, t := range targets {
		for k := range t.Env {
			keySet[k] = struct{}{}
		}
	}

	names := make([]string, len(targets))
	for i, t := range targets {
		names[i] = t.Name
	}

	var entries []GradientEntry
	for k := range keySet {
		values := make([]string, len(targets))
		for i, t := range targets {
			values[i] = t.Env[k]
		}

		changes := 0
		for i := 1; i < len(values); i++ {
			if values[i] != values[i-1] {
				changes++
			}
		}

		if !opts.IncludeStable && changes < opts.MinChanges {
			continue
		}
		if opts.IncludeStable && changes < opts.MinChanges && changes > 0 {
			continue
		}

		entries = append(entries, GradientEntry{
			Key:       k,
			Steps:     names,
			Values:    values,
			Changes:   changes,
			Direction: gradientDirection(values),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Changes != entries[j].Changes {
			return entries[i].Changes > entries[j].Changes
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// HasGradientChanges returns true when at least one entry has changes.
func HasGradientChanges(entries []GradientEntry) bool {
	for _, e := range entries {
		if e.Changes > 0 {
			return true
		}
	}
	return false
}

func gradientDirection(values []string) string {
	if len(values) < 2 {
		return "stable"
	}
	changes := 0
	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1] {
			changes++
		}
	}
	if changes == 0 {
		return "stable"
	}
	if changes == len(values)-1 {
		return "volatile"
	}
	return "shifting"
}

// NamedEnv pairs a target name with its env map.
type NamedEnv struct {
	Name string
	Env  map[string]string
}
