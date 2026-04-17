package differ

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ReportOptions controls output format for diff reports.
type ReportOptions struct {
	Format    string // "text" or "json"
	ShowOnly  string // "all", "changed", "added", "removed"
}

// DefaultReportOptions returns sensible defaults.
func DefaultReportOptions() ReportOptions {
	return ReportOptions{Format: "text", ShowOnly: "all"}
}

// ReportDiff writes a human-readable or JSON diff report to w.
func ReportDiff(w io.Writer, results []Result, opts ReportOptions) error {
	filtered := filterResults(results, opts.ShowOnly)
	switch opts.Format {
	case "json":
		return writeDiffJSON(w, filtered)
	default:
		return writeDiffText(w, filtered)
	}
}

func filterResults(results []Result, showOnly string) []Result {
	if showOnly == "all" || showOnly == "" {
		return results
	}
	var out []Result
	for _, r := range results {
		switch showOnly {
		case "changed":
			if r.Status == StatusChanged {
				out = append(out, r)
			}
		case "added":
			if r.Status == StatusAdded {
				out = append(out, r)
			}
		case "removed":
			if r.Status == StatusRemoved {
				out = append(out, r)
			}
		}
	}
	return out
}

func writeDiffText(w io.Writer, results []Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no differences found")
		return err
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Key < results[j].Key })
	for _, r := range results {
		switch r.Status {
		case StatusAdded:
			fmt.Fprintf(w, "+ %s=%s\n", r.Key, r.NewValue)
		case StatusRemoved:
			fmt.Fprintf(w, "- %s=%s\n", r.Key, r.OldValue)
		case StatusChanged:
			fmt.Fprintf(w, "~ %s: %s -> %s\n", r.Key, r.OldValue, r.NewValue)
		case StatusUnchanged:
			fmt.Fprintf(w, "  %s=%s\n", r.Key, r.NewValue)
		}
	}
	return nil
}

func writeDiffJSON(w io.Writer, results []Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if results == nil {
		results = []Result{}
	}
	return enc.Encode(results)
}
