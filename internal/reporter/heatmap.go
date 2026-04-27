package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// HeatmapReportOptions controls how heatmap results are rendered.
type HeatmapReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultHeatmapReportOptions returns sensible defaults.
func DefaultHeatmapReportOptions(w io.Writer) HeatmapReportOptions {
	return HeatmapReportOptions{Format: "text", Writer: w}
}

// ReportHeatmap writes heatmap entries to the configured writer.
func ReportHeatmap(entries []differ.HeatmapEntry, opts HeatmapReportOptions) error {
	switch opts.Format {
	case "json":
		return writeHeatmapJSON(entries, opts.Writer)
	default:
		return writeHeatmapText(entries, opts.Writer)
	}
}

func writeHeatmapText(entries []differ.HeatmapEntry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "no heatmap entries")
		return err
	}
	for _, e := range entries {
		targets := strings.Join(e.Targets, ", ")
		_, err := fmt.Fprintf(w, "%-40s changes=%-4d targets=[%s]\n", e.Key, e.Changes, targets)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeHeatmapJSON(entries []differ.HeatmapEntry, w io.Writer) error {
	type row struct {
		Key     string   `json:"key"`
		Changes int      `json:"changes"`
		Targets []string `json:"targets"`
	}
	rows := make([]row, len(entries))
	for i, e := range entries {
		rows[i] = row{Key: e.Key, Changes: e.Changes, Targets: e.Targets}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
