package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envlens/internal/differ"
)

// MultiDiffOptions controls output of ReportMultiDiff.
type MultiDiffOptions struct {
	Format     string // "text" or "json"
	ShowAll    bool   // include unchanged entries
}

// DefaultMultiDiffOptions returns sensible defaults.
func DefaultMultiDiffOptions() MultiDiffOptions {
	return MultiDiffOptions{Format: "text", ShowAll: false}
}

// ReportMultiDiff writes a multi-target diff report to w.
func ReportMultiDiff(w io.Writer, diffs []differ.TargetDiff, opts MultiDiffOptions) error {
	switch opts.Format {
	case "json":
		return writeMultiDiffJSON(w, diffs, opts)
	default:
		return writeMultiDiffText(w, diffs, opts)
	}
}

func writeMultiDiffText(w io.Writer, diffs []differ.TargetDiff, opts MultiDiffOptions) error {
	for _, td := range diffs {
		fmt.Fprintf(w, "=== %s ===\n", td.Target)
		printed := 0
		for _, r := range td.Results {
			if !opts.ShowAll && r.Status == differ.StatusUnchanged {
				continue
			}
			var prefix string
			switch r.Status {
			case differ.StatusAdded:
				prefix = "+"
			case differ.StatusRemoved:
				prefix = "-"
			case differ.StatusChanged:
				prefix = "~"
			default:
				prefix = " "
			}
			fmt.Fprintf(w, "  %s %s: %q -> %q\n", prefix, r.Key, r.OldValue, r.NewValue)
			printed++
		}
		if printed == 0 {
			fmt.Fprintf(w, "  (no changes)\n")
		}
	}
	return nil
}

func writeMultiDiffJSON(w io.Writer, diffs []differ.TargetDiff, opts MultiDiffOptions) error {
	type entry struct {
		Key      string `json:"key"`
		Status   string `json:"status"`
		OldValue string `json:"old_value"`
		NewValue string `json:"new_value"`
	}
	type targetOut struct {
		Target  string  `json:"target"`
		Entries []entry `json:"entries"`
	}
	var out []targetOut
	for _, td := range diffs {
		var entries []entry
		for _, r := range td.Results {
			if !opts.ShowAll && r.Status == differ.StatusUnchanged {
				continue
			}
			entries = append(entries, entry{Key: r.Key, Status: string(r.Status), OldValue: r.OldValue, NewValue: r.NewValue})
		}
		out = append(out, targetOut{Target: td.Target, Entries: entries})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
