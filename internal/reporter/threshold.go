package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envlens/internal/differ"
)

// ThresholdReportOptions controls output for threshold results.
type ThresholdReportOptions struct {
	Format string // "text" or "json"
	Out    io.Writer
}

// DefaultThresholdReportOptions returns text output to stdout.
func DefaultThresholdReportOptions(w io.Writer) ThresholdReportOptions {
	return ThresholdReportOptions{Format: "text", Out: w}
}

// ReportThreshold writes threshold-filtered diff results.
func ReportThreshold(results []differ.ThresholdResult, opts ThresholdReportOptions) error {
	switch opts.Format {
	case "json":
		return writeThresholdJSON(results, opts.Out)
	default:
		return writeThresholdText(results, opts.Out)
	}
}

func writeThresholdText(results []differ.ThresholdResult, w io.Writer) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no targets exceeded change threshold")
		return err
	}
	for _, tr := range results {
		fmt.Fprintf(w, "[%s] %d change(s)\n", tr.Target, tr.Count)
		for _, r := range tr.Results {
			switch r.Status {
			case "added":
				fmt.Fprintf(w, "  + %s=%s\n", r.Key, r.NewValue)
			case "removed":
				fmt.Fprintf(w, "  - %s\n", r.Key)
			case "changed":
				fmt.Fprintf(w, "  ~ %s: %s -> %s\n", r.Key, r.OldValue, r.NewValue)
			default:
				fmt.Fprintf(w, "    %s\n", r.Key)
			}
		}
	}
	return nil
}

func writeThresholdJSON(results []differ.ThresholdResult, w io.Writer) error {
	type jsonEntry struct {
		Target  string          `json:"target"`
		Count   int             `json:"count"`
		Results []differ.Result `json:"results"`
	}
	out := make([]jsonEntry, 0, len(results))
	for _, tr := range results {
		out = append(out, jsonEntry{Target: tr.Target, Count: tr.Count, Results: tr.Results})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
