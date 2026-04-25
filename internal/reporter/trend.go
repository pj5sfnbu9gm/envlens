package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// TrendReportOptions controls how trend reports are rendered.
type TrendReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultTrendReportOptions returns sensible defaults.
func DefaultTrendReportOptions(w io.Writer) TrendReportOptions {
	return TrendReportOptions{Format: "text", Writer: w}
}

// ReportTrend writes a trend report to the configured writer.
func ReportTrend(entries []differ.TrendEntry, opts TrendReportOptions) error {
	switch strings.ToLower(opts.Format) {
	case "json":
		return writeTrendJSON(entries, opts.Writer)
	default:
		return writeTrendText(entries, opts.Writer)
	}
}

func writeTrendText(entries []differ.TrendEntry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "no trend data")
		return err
	}
	for _, e := range entries {
		countStrs := make([]string, len(e.Counts))
		for i, c := range e.Counts {
			countStrs[i] = fmt.Sprintf("%d", c)
		}
		_, err := fmt.Fprintf(w, "%-40s total=%-4d direction=%-6s windows=[%s]\n",
			e.Key, e.Total, e.Direction, strings.Join(countStrs, ","))
		if err != nil {
			return err
		}
	}
	return nil
}

func writeTrendJSON(entries []differ.TrendEntry, w io.Writer) error {
	type jsonEntry struct {
		Key       string                 `json:"key"`
		Counts    []int                  `json:"counts"`
		Total     int                    `json:"total"`
		Direction differ.TrendDirection  `json:"direction"`
	}
	out := make([]jsonEntry, len(entries))
	for i, e := range entries {
		out[i] = jsonEntry{Key: e.Key, Counts: e.Counts, Total: e.Total, Direction: e.Direction}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
