package differ

// LabelEntry associates a diff result with a human-readable label and optional metadata.
type LabelEntry struct {
	Key    string
	Status string
	Old    string
	New    string
	Label  string
	Meta   map[string]string
}

// LabelOptions controls how labels are applied to diff results.
type LabelOptions struct {
	// Labels maps exact keys to custom label strings.
	Labels map[string]string
	// DefaultLabel is used when no specific label is found.
	DefaultLabel string
	// IncludeUnchanged includes unchanged results in the output.
	IncludeUnchanged bool
}

// DefaultLabelOptions returns sensible defaults for LabelOptions.
func DefaultLabelOptions() LabelOptions {
	return LabelOptions{
		Labels:           map[string]string{},
		DefaultLabel:     "unlabeled",
		IncludeUnchanged: false,
	}
}

// Label applies human-readable labels to a slice of DiffResults.
// Keys present in opts.Labels receive that label; all others get opts.DefaultLabel.
func Label(results []Result, opts LabelOptions) []LabelEntry {
	out := make([]LabelEntry, 0, len(results))
	for _, r := range results {
		if r.Status == StatusUnchanged && !opts.IncludeUnchanged {
			continue
		}
		lbl, ok := opts.Labels[r.Key]
		if !ok {
			lbl = opts.DefaultLabel
		}
		out = append(out, LabelEntry{
			Key:    r.Key,
			Status: string(r.Status),
			Old:    r.OldValue,
			New:    r.NewValue,
			Label:  lbl,
			Meta:   map[string]string{},
		})
	}
	return out
}

// HasLabeledChanges returns true if any LabelEntry reflects a non-unchanged status.
func HasLabeledChanges(entries []LabelEntry) bool {
	for _, e := range entries {
		if e.Status != string(StatusUnchanged) {
			return true
		}
	}
	return false
}
