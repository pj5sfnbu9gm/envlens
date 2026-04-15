// Package reporter formats and outputs diff results for environment
// variable comparisons across deployment targets.
package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Options configures the reporter output.
type Options struct {
	Format  Format
	NoColor bool
	Writer  io.Writer
}

// DefaultOptions returns sensible defaults for reporter options.
func DefaultOptions() Options {
	return Options{
		Format:  FormatText,
		NoColor: false,
		Writer:  os.Stdout,
	}
}

// Report writes a formatted diff result to the configured writer.
func Report(results []differ.Result, from, to string, opts Options) error {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	switch opts.Format {
	case FormatJSON:
		return writeJSON(w, results, from, to)
	default:
		return writeText(w, results, from, to, !opts.NoColor)
	}
}

func writeText(w io.Writer, results []differ.Result, from, to string, color bool) error {
	fmt.Fprintf(w, "envlens diff: %s → %s\n", from, to)
	fmt.Fprintln(w, strings.Repeat("─", 40))

	changed := 0
	for _, r := range results {
		if r.Status == differ.StatusUnchanged {
			continue
		}
		changed++
		line := formatTextLine(r, color)
		fmt.Fprintln(w, line)
	}

	if changed == 0 {
		fmt.Fprintln(w, "No differences found.")
	}
	return nil
}

func formatTextLine(r differ.Result, color bool) string {
	switch r.Status {
	case differ.StatusAdded:
		if color {
			return fmt.Sprintf("\033[32m+ %s=%s\033[0m", r.Key, r.ToValue)
		}
		return fmt.Sprintf("+ %s=%s", r.Key, r.ToValue)
	case differ.StatusRemoved:
		if color {
			return fmt.Sprintf("\033[31m- %s=%s\033[0m", r.Key, r.FromValue)
		}
		return fmt.Sprintf("- %s=%s", r.Key, r.FromValue)
	case differ.StatusChanged:
		if color {
			return fmt.Sprintf("\033[33m~ %s: %s → %s\033[0m", r.Key, r.FromValue, r.ToValue)
		}
		return fmt.Sprintf("~ %s: %s → %s", r.Key, r.FromValue, r.ToValue)
	}
	return ""
}
