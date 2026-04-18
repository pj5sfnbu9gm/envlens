package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/user/envlens/internal/differ"
)

// ScopeReportOptions configures the scope diff reporter.
type ScopeReportOptions struct {
	Format string // "text" or "json"
	Out    io.Writer
}

// DefaultScopeOptions returns sensible defaults for scope reporting.
func DefaultScopeOptions() ScopeReportOptions {
	return ScopeReportOptions{Format: "text", Out: os.Stdout}
}

// ReportScope writes a human-readable or JSON report of scoped diff results.
func ReportScope(results []differ.ScopeResult, opts ScopeReportOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	switch opts.Format {
	case "json":
		return writeScopeJSON(opts.Out, results)
	default:
		return writeScopeText(opts.Out, results)
	}
}

func writeScopeText(w io.Writer, results []differ.ScopeResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no scoped results")
		return err
	}
	for _, r := range results {
		var line string
		switch r.Status {
		case "added":
			line = fmt.Sprintf("+ %s=%s", r.Key, r.New)
		case "removed":
			line = fmt.Sprintf("- %s=%s", r.Key, r.Old)
		case "changed":
			line = fmt.Sprintf("~ %s: %s -> %s", r.Key, r.Old, r.New)
		default:
			line = fmt.Sprintf("  %s=%s", r.Key, r.New)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeScopeJSON(w io.Writer, results []differ.ScopeResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
