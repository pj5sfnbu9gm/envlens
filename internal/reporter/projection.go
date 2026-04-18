package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envlens/internal/differ"
)

// ProjectionReportOptions controls output of ReportProjection.
type ProjectionReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultProjectionOptions returns sensible defaults for projection reporting.
func DefaultProjectionReportOptions(w io.Writer) ProjectionReportOptions {
	return ProjectionReportOptions{Format: "text", Writer: w}
}

// ReportProjection writes projected diff results to the configured writer.
func ReportProjection(results []differ.DiffResult, opts ProjectionReportOptions) error {
	switch opts.Format {
	case "json":
		return writeProjectionJSON(results, opts.Writer)
	default:
		return writeProjectionText(results, opts.Writer)
	}
}

func writeProjectionText(results []differ.DiffResult, w io.Writer) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no projected results")
		return err
	}
	for _, r := range results {
		var line string
		switch r.Status {
		case differ.StatusAdded:
			line = fmt.Sprintf("+ %s = %s", r.Key, r.NewValue)
		case differ.StatusRemoved:
			line = fmt.Sprintf("- %s = %s", r.Key, r.OldValue)
		case differ.StatusChanged:
			line = fmt.Sprintf("~ %s: %s -> %s", r.Key, r.OldValue, r.NewValue)
		default:
			line = fmt.Sprintf("  %s = %s", r.Key, r.NewValue)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeProjectionJSON(results []differ.DiffResult, w io.Writer) error {
	type entry struct {
		Key      string `json:"key"`
		Status   string `json:"status"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
	}
	out := make([]entry, len(results))
	for i, r := range results {
		out[i] = entry{Key: r.Key, Status: string(r.Status), OldValue: r.OldValue, NewValue: r.NewValue}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
