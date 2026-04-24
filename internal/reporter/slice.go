package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/envlens/internal/differ"
)

// SliceReportOptions configures ReportSlice output.
type SliceReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultSliceReportOptions returns text output to stdout.
func DefaultSliceReportOptions() SliceReportOptions {
	return SliceReportOptions{
		Format: "text",
		Writer: os.Stdout,
	}
}

// ReportSlice writes grouped slice diff results to the configured writer.
func ReportSlice(slices []differ.SliceResult, opts SliceReportOptions) error {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}
	switch opts.Format {
	case "json":
		return writeSliceJSON(w, slices)
	default:
		return writeSliceText(w, slices)
	}
}

func writeSliceText(w io.Writer, slices []differ.SliceResult) error {
	if len(slices) == 0 {
		_, err := fmt.Fprintln(w, "no slice changes")
		return err
	}
	for _, s := range slices {
		fmt.Fprintf(w, "[%s]\n", s.Prefix)
		for _, r := range s.Results {
			switch r.Status {
			case differ.StatusAdded:
				fmt.Fprintf(w, "  + %s = %s\n", r.Key, r.NewValue)
			case differ.StatusRemoved:
				fmt.Fprintf(w, "  - %s = %s\n", r.Key, r.OldValue)
			case differ.StatusChanged:
				fmt.Fprintf(w, "  ~ %s: %s -> %s\n", r.Key, r.OldValue, r.NewValue)
			default:
				fmt.Fprintf(w, "    %s = %s\n", r.Key, r.NewValue)
			}
		}
	}
	return nil
}

func writeSliceJSON(w io.Writer, slices []differ.SliceResult) error {
	type entry struct {
		Prefix  string         `json:"prefix"`
		Results []differ.Result `json:"results"`
	}
	out := make([]entry, 0, len(slices))
	for _, s := range slices {
		out = append(out, entry{Prefix: s.Prefix, Results: s.Results})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
