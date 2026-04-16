// Package pivotter transposes a multi-target environment comparison into
// a key-centric view: for each env key, show its value across all targets.
package pivotter

import "sort"

// Row holds the value of a single key across every target.
type Row struct {
	Key     string
	Targets map[string]string // target name -> value ("" if absent)
}

// Options controls Pivot behaviour.
type Options struct {
	// IncludeUnchanged includes keys whose value is identical in all targets.
	IncludeUnchanged bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{IncludeUnchanged: false}
}

// Pivot takes a map of target->envMap and returns one Row per key.
// When IncludeUnchanged is false, keys that share the same value in
// every target are omitted.
func Pivot(targets map[string]map[string]string, opts Options) []Row {
	// Collect all keys.
	keySet := map[string]struct{}{}
	for _, env := range targets {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	// Build target name list for stable ordering.
	targetNames := make([]string, 0, len(targets))
	for name := range targets {
		targetNames = append(targetNames, name)
	}
	sort.Strings(targetNames)

	var rows []Row
	for key := range keySet {
		row := Row{Key: key, Targets: make(map[string]string, len(targets))}
		var first string
		allSame := true
		for i, name := range targetNames {
			v := targets[name][key]
			row.Targets[name] = v
			if i == 0 {
				first = v
			} else if v != first {
				allSame = false
			}
		}
		if !allSame || opts.IncludeUnchanged {
			rows = append(rows, row)
		}
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].Key < rows[j].Key })
	return rows
}
