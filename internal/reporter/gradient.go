package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// GradientReportOptions controls gradient report output.
type GradientReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultGradientReportOptions returns text output to stdout.
func DefaultGradientReportOptions(w io.Writer) GradientReportOptions {
	return GradientReportOptions{Format: "text", Writer: w}
}

// ReportGradient writes a gradient analysis report.
func ReportGradient(entries []differ.GradientEntry, opts GradientReportOptions) error {
	switch strings.ToLower(opts.Format) {
	case "json":
		return writeGradientJSON(entries, opts.Writer)
	default:
		return writeGradientText(entries, opts.Writer)
	}
}

func writeGradientText(entries []differ.GradientEntry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "gradient: no shifting keys detected")
		return err
	}
	for _, e := range entries {
		_, err := fmt.Fprintf(w, "[%s] %s  changes=%d\n", e.Direction, e.Key, e.Changes)
		if err != nil {
			return err
		}
		for i, step := range e.Steps {
			val := e.Values[i]
			if val == "" {
				val = "<absent>"
			}
			_, err = fmt.Fprintf(w, "  %-12s %s\n", step+":", val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeGradientJSON(entries []differ.GradientEntry, w io.Writer) error {
	type row struct {
		Key       string            `json:"key"`
		Changes   int               `json:"changes"`
		Direction string            `json:"direction"`
		Values    map[string]string `json:"values"`
	}
	rows := make([]row, len(entries))
	for i, e := range entries {
		vals := make(map[string]string, len(e.Steps))
		for j, s := range e.Steps {
			vals[s] = e.Values[j]
		}
		rows[i] = row{Key: e.Key, Changes: e.Changes, Direction: e.Direction, Values: vals}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
