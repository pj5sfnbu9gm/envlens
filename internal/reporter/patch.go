package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourorg/envlens/internal/patcher"
)

// PatchOptions controls how patch results are rendered.
type PatchOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultPatchOptions returns sensible defaults.
func DefaultPatchOptions(w io.Writer) PatchOptions {
	return PatchOptions{Format: "text", Writer: w}
}

// ReportPatch writes patch results to the configured writer.
func ReportPatch(results []patcher.Result, opts PatchOptions) error {
	switch opts.Format {
	case "json":
		return writePatchJSON(results, opts.Writer)
	default:
		return writePatchText(results, opts.Writer)
	}
}

func writePatchText(results []patcher.Result, w io.Writer) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no patch operations")
		return err
	}
	for _, r := range results {
		status := "skipped"
		if r.Applied {
			status = "applied"
		}
		line := fmt.Sprintf("[%s] %-8s %s", status, r.Patch.Op, r.Patch.Key)
		if r.Note != "" {
			line += fmt.Sprintf(" (%s)", r.Note)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writePatchJSON(results []patcher.Result, w io.Writer) error {
	type row struct {
		Op      string `json:"op"`
		Key     string `json:"key"`
		To      string `json:"to,omitempty"`
		Value   string `json:"value,omitempty"`
		Applied bool   `json:"applied"`
		Note    string `json:"note,omitempty"`
	}
	rows := make([]row, len(results))
	for i, r := range results {
		rows[i] = row{
			Op:      string(r.Patch.Op),
			Key:     r.Patch.Key,
			To:      r.Patch.To,
			Value:   r.Patch.Value,
			Applied: r.Applied,
			Note:    r.Note,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
