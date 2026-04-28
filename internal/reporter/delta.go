package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envlens/internal/differ"
)

// DeltaReportOptions controls how delta results are rendered.
type DeltaReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultDeltaReportOptions returns text-format defaults writing to stdout.
func DefaultDeltaReportOptions(w io.Writer) DeltaReportOptions {
	return DeltaReportOptions{Format: "text", Writer: w}
}

// ReportDelta writes a formatted delta report to the configured writer.
func ReportDelta(entries []differ.DeltaEntry, opts DeltaReportOptions) error {
	switch opts.Format {
	case "json":
		return writeDeltaJSON(entries, opts.Writer)
	default:
		return writeDeltaText(entries, opts.Writer)
	}
}

func writeDeltaText(entries []differ.DeltaEntry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "no delta entries")
		return err
	}
	_, err := fmt.Fprintf(w, "%-30s %6s %8s %8s %6s\n", "KEY", "ADDS", "REMOVALS", "CHANGES", "NET")
	if err != nil {
		return err
	}
	for _, e := range entries {
		_, err = fmt.Fprintf(w, "%-30s %6d %8d %8d %+6d\n",
			e.Key, e.Adds, e.Removals, e.Changes, e.Net)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDeltaJSON(entries []differ.DeltaEntry, w io.Writer) error {
	type row struct {
		Key      string `json:"key"`
		Adds     int    `json:"adds"`
		Removals int    `json:"removals"`
		Changes  int    `json:"changes"`
		Net      int    `json:"net"`
	}
	rows := make([]row, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, row{Key: e.Key, Adds: e.Adds, Removals: e.Removals, Changes: e.Changes, Net: e.Net})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
