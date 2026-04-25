package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envlens/internal/differ"
)

// SignalReportOptions controls how signal results are rendered.
type SignalReportOptions struct {
	Format  string // "text" or "json"
	Writer  io.Writer
}

// DefaultSignalReportOptions returns text-format defaults writing to stdout.
func DefaultSignalReportOptions(w io.Writer) SignalReportOptions {
	return SignalReportOptions{Format: "text", Writer: w}
}

// ReportSignal renders signal entries to the configured output.
func ReportSignal(entries []differ.SignalEntry, opts SignalReportOptions) error {
	switch strings.ToLower(opts.Format) {
	case "json":
		return writeSignalJSON(entries, opts.Writer)
	default:
		return writeSignalText(entries, opts.Writer)
	}
}

func writeSignalText(entries []differ.SignalEntry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "no high-signal keys detected")
		return err
	}
	for _, e := range entries {
		line := fmt.Sprintf("[SIGNAL] %s  changes=%d  targets=[%s]",
			e.Key, e.ChangeCount, strings.Join(e.Targets, ","))
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeSignalJSON(entries []differ.SignalEntry, w io.Writer) error {
	out := entries
	if out == nil {
		out = []differ.SignalEntry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
