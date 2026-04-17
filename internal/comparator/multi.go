package comparator

import (
	"fmt"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/loader"
)

// TargetResult holds the diff results for a single named target.
type TargetResult struct {
	Name    string
	Results []differ.Result
}

// MultiCompareOptions configures MultiCompareAll behavior.
type MultiCompareOptions struct {
	Baseline string
	Targets  map[string]string // name -> file path
}

// MultiCompareAll loads the baseline and all targets, then diffs each target
// against the baseline. Returns one TargetResult per target.
func MultiCompareAll(opts MultiCompareOptions) ([]TargetResult, error) {
	if opts.Baseline == "" {
		return nil, fmt.Errorf("comparator: baseline path is required")
	}
	if len(opts.Targets) == 0 {
		return nil, fmt.Errorf("comparator: at least one target is required")
	}

	base, err := loader.LoadFile(opts.Baseline)
	if err != nil {
		return nil, fmt.Errorf("comparator: loading baseline %q: %w", opts.Baseline, err)
	}

	var out []TargetResult
	for name, path := range opts.Targets {
		target, err := loader.LoadFile(path)
		if err != nil {
			return nil, fmt.Errorf("comparator: loading target %q (%s): %w", name, path, err)
		}
		results := differ.Diff(base, target)
		out = append(out, TargetResult{Name: name, Results: results})
	}
	return out, nil
}

// AnyTargetHasChanges returns true if at least one TargetResult contains
// a non-Unchanged diff entry.
func AnyTargetHasChanges(results []TargetResult) bool {
	for _, tr := range results {
		if HasChanges(tr.Results) {
			return true
		}
	}
	return false
}
