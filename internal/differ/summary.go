package differ

// Summary holds aggregate counts from a diff result set.
type Summary struct {
	Total     int `json:"total"`
	Added     int `json:"added"`
	Removed   int `json:"removed"`
	Changed   int `json:"changed"`
	Unchanged int `json:"unchanged"`
}

// Summarize computes aggregate statistics from a slice of Results.
func Summarize(results []Result) Summary {
	var s Summary
	for _, r := range results {
		s.Total++
		switch r.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusChanged:
			s.Changed++
		case StatusUnchanged:
			s.Unchanged++
		}
	}
	return s
}

// HasDifferences returns true if any results are not unchanged.
func (s Summary) HasDifferences() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}
