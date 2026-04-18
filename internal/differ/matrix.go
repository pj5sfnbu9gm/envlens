package differ

// MatrixEntry holds the diff result between two named targets.
type MatrixEntry struct {
	From    string
	To      string
	Results []Result
}

// Matrix computes pairwise diffs between all provided targets.
// Each unique ordered pair (A→B) is compared using Diff.
func Matrix(targets map[string]map[string]string) []MatrixEntry {
	names := sortStrings(keys(targets))
	var entries []MatrixEntry
	for i := 0; i < len(names); i++ {
		for j := 0; j < len(names); j++ {
			if i == j {
				continue
			}
			from := names[i]
			to := names[j]
			results := Diff(targets[from], targets[to])
			entries = append(entries, MatrixEntry{
				From:    from,
				To:      to,
				Results: results,
			})
		}
	}
	return entries
}

// HasMatrixChanges returns true if any pair in the matrix has non-unchanged results.
func HasMatrixChanges(entries []MatrixEntry) bool {
	for _, e := range entries {
		for _, r := range e.Results {
			if r.Status != StatusUnchanged {
				return true
			}
		}
	}
	return false
}

func keys(m map[string]map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
