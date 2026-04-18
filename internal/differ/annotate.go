package differ

// Annotation holds a label attached to a diff result key.
type Annotation struct {
	Key   string
	Label string
	Note  string
}

// AnnotateOptions controls how annotations are applied.
type AnnotateOptions struct {
	// Rules maps a status string ("added", "removed", "changed", "unchanged")
	// to a label.
	StatusLabels map[string]string
	// KeyNotes maps specific keys to free-form notes.
	KeyNotes map[string]string
}

// DefaultAnnotateOptions returns sensible defaults.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		StatusLabels: map[string]string{
			"added":     "[+]",
			"removed":   "[-]",
			"changed":   "[~]",
			"unchanged": "[=]",
		},
		KeyNotes: map[string]string{},
	}
}

// Annotate attaches labels and notes to a slice of Result values.
func Annotate(results []Result, opts AnnotateOptions) []Annotation {
	annotations := make([]Annotation, 0, len(results))
	for _, r := range results {
		label := opts.StatusLabels[string(r.Status)]
		note := opts.KeyNotes[r.Key]
		annotations = append(annotations, Annotation{
			Key:   r.Key,
			Label: label,
			Note:  note,
		})
	}
	return annotations
}
