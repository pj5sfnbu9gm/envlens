package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/user/envlens/internal/differ"
)

// RollupOptions controls rollup report output.
type RollupOptions struct {
	Format string // "text" or "json"
	Out    io.Writer
}

// DefaultRollupOptions returns sensible defaults.
func DefaultRollupOptions() RollupOptions {
	return RollupOptions{Format: "text", Out: os.Stdout}
}

// ReportRollup writes a rollup summary to the configured output.
func ReportRollup(entries []differ.RollupEntry, opts RollupOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	switch opts.Format {
	case "json":
		return writeRollupJSON(entries, opts.Out)
	default:
		return writeRollupText(entries, opts.Out)
	}
}

func writeRollupText(entries []differ.RollupEntry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "no rollup changes")
		return err
	}
	fmt.Fprintf(w, "%-30s %7s %7s %7s %7s\n", "KEY", "ADDED", "REMOVED", "CHANGED", "TOTAL")
	for _, e := range entries {
		_, err := fmt.Fprintf(w, "%-30s %7d %7d %7d %7d\n",
			e.Key, e.Added, e.Removed, e.Changed, e.Total)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeRollupJSON(entries []differ.RollupEntry, w io.Writer) error {
	type row struct {
		Key       string `json:"key"`
		Added     int    `json:"added"`
		Removed   int    `json:"removed"`
		Changed   int    `json:"changed"`
		Unchanged int    `json:"unchanged"`
		Total     int    `json:"total"`
	}
	rows := make([]row, len(entries))
	for i, e := range entries {
		rows[i] = row{e.Key, e.Added, e.Removed, e.Changed, e.Unchanged, e.Total}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
