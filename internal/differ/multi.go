package differ

import "fmt"

// TargetDiff holds the diff results for a single named target against a baseline.
type TargetDiff struct {
	Target  string
	Results []Result
}

// MultiDiff compares multiple target env maps against a single baseline.
// targets is a map of target-name -> env map.
func MultiDiff(baseline map[string]string, targets map[string]map[string]string) ([]TargetDiff, error) {
	if baseline == nil {
		return nil, fmt.Errorf("baseline must not be nil")
	}
	var out []TargetDiff
	// stable order
	names := make([]string, 0, len(targets))
	for name := range targets {
		names = append(names, name)
	}
	sortStrings(names)
	for _, name := range names {
		results := Diff(baseline, targets[name])
		out = append(out, TargetDiff{
			Target:  name,
			Results: results,
		})
	}
	return out, nil
}

// HasAnyChanges returns true if any target diff contains a non-unchanged result.
func HasAnyChanges(diffs []TargetDiff) bool {
	for _, td := range diffs {
		for _, r := range td.Results {
			if r.Status != StatusUnchanged {
				return true
			}
		}
	}
	return false
}
