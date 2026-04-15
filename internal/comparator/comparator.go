package comparator

import (
	"fmt"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/loader"
	"github.com/user/envlens/internal/resolver"
)

// TargetDiff holds the diff results between two named targets.
type TargetDiff struct {
	From    string
	To      string
	Results []differ.Result
}

// CompareAll loads all resolved targets and produces pairwise diffs between
// a baseline target and every other target.
func CompareAll(baseline string, targets []resolver.Target) ([]TargetDiff, error) {
	envMaps := make(map[string]map[string]string, len(targets))
	for _, t := range targets {
		m, err := loader.LoadFile(t.Path)
		if err != nil {
			return nil, fmt.Errorf("comparator: loading target %q: %w", t.Name, err)
		}
		envMaps[t.Name] = m
	}

	baseMap, ok := envMaps[baseline]
	if !ok {
		return nil, fmt.Errorf("comparator: baseline target %q not found", baseline)
	}

	var diffs []TargetDiff
	for _, t := range targets {
		if t.Name == baseline {
			continue
		}
		results := differ.Diff(baseMap, envMaps[t.Name])
		diffs = append(diffs, TargetDiff{
			From:    baseline,
			To:      t.Name,
			Results: results,
		})
	}
	return diffs, nil
}

// HasChanges returns true if any TargetDiff contains non-unchanged results.
func HasChanges(diffs []TargetDiff) bool {
	for _, d := range diffs {
		for _, r := range d.Results {
			if r.Status != differ.Unchanged {
				return true
			}
		}
	}
	return false
}
