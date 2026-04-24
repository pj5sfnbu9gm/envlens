package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/user/envlens/internal/differ"
)

// LabelReportOptions controls output of the label report.
type LabelReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultLabelReportOptions returns defaults writing text to stdout.
func DefaultLabelReportOptions() LabelReportOptions {
	return LabelReportOptions{
		Format: "text",
		Writer: os.Stdout,
	}
}

// ReportLabel writes labeled diff entries in the chosen format.
func ReportLabel(entries []differ.LabelEntry, opts LabelReportOptions) error {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}
	switch opts.Format {
	case "json":
		return writeLabelJSON(w, entries)
	default:
		return writeLabelText(w, entries)
	}
}

func writeLabelText(w io.Writer, entries []differ.LabelEntry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "no labeled changes")
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tSTATUS\tLABEL\tOLD\tNEW")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			e.Key, e.Status, e.Label, e.Old, e.New)
	}
	return tw.Flush()
}

func writeLabelJSON(w io.Writer, entries []differ.LabelEntry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	type row struct {
		Key    string `json:"key"`
		Status string `json:"status"`
		Label  string `json:"label"`
		Old    string `json:"old,omitempty"`
		New    string `json:"new,omitempty"`
	}
	rows := make([]row, len(entries))
	for i, e := range entries {
		rows[i] = row{Key: e.Key, Status: e.Status, Label: e.Label, Old: e.Old, New: e.New}
	}
	return enc.Encode(rows)
}
