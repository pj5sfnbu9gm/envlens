package differ

import "sort"

// WindowOptions controls how a sliding window diff is applied.
type WindowOptions struct {
	// Size is the number of consecutive snapshots to include in each window.
	// Must be >= 2.
	Size int

	// IgnoreUnchanged skips keys that are identical across the entire window.
	IgnoreUnchanged bool
}

// WindowResult holds the diff results for a single window of snapshots.
type WindowResult struct {
	// Labels are the names of the snapshots in this window (e.g. ["t0", "t1"]).
	Labels []string
	// Results contains per-key diff entries for this window.
	Results []Result
}

// DefaultWindowOptions returns sensible defaults for Window.
func DefaultWindowOptions() WindowOptions {
	return WindowOptions{
		Size:            2,
		IgnoreUnchanged: true,
	}
}

// Window applies a sliding window over an ordered list of named env maps,
// returning one WindowResult per consecutive pair (or group of Size).
// Snapshots must be provided in chronological order.
func Window(snapshots []NamedEnv, opts WindowOptions) []WindowResult {
	if opts.Size < 2 {
		opts.Size = 2
	}
	if len(snapshots) < opts.Size {
		return nil
	}

	var out []WindowResult
	for i := 0; i <= len(snapshots)-opts.Size; i++ {
		window := snapshots[i : i+opts.Size]
		base := window[0].Env
		target := window[opts.Size-1].Env

		results := Diff(base, target)
		if opts.IgnoreUnchanged {
			filtered := results[:0]
			for _, r := range results {
				if r.Status != StatusUnchanged {
					filtered = append(filtered, r)
				}
			}
			results = filtered
		}

		labels := make([]string, opts.Size)
		for j, s := range window {
			labels[j] = s.Name
		}

		out = append(out, WindowResult{
			Labels:  labels,
			Results: results,
		})
	}
	return out
}

// HasWindowChanges returns true if any window contains at least one non-unchanged result.
func HasWindowChanges(windows []WindowResult) bool {
	for _, w := range windows {
		for _, r := range w.Results {
			if r.Status != StatusUnchanged {
				return true
			}
		}
	}
	return false
}

// NamedEnv pairs a label with an environment map.
type NamedEnv struct {
	Name string
	Env  map[string]string
}

// ensure sort is used (stability)
var _ = sort.Strings
