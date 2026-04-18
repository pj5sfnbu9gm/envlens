package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envlens/internal/differ"
)

// MatrixOptions controls output of ReportMatrix.
type MatrixOptions struct {
	Format string // "text" or "json"
	ShowUnchanged bool
}

// DefaultMatrixOptions returns sensible defaults.
func DefaultMatrixOptions() MatrixOptions {
	return MatrixOptions{Format: "text", ShowUnchanged: false}
}

// ReportMatrix writes a pairwise diff matrix to w.
func ReportMatrix(w io.Writer, entries []differ.MatrixEntry, opts MatrixOptions) error {
	if opts.Format == "json" {
		return writeMatrixJSON(w, entries, opts)
	}
	return writeMatrixText(w, entries, opts)
}

func writeMatrixText(w io.Writer, entries []differ.MatrixEntry, opts MatrixOptions) error {
	for _, e := range entries {
		fmt.Fprintf(w, "[%s → %s]\n", e.From, e.To)
		printed := 0
		for _, r := range e.Results {
			if !opts.ShowUnchanged && r.Status == differ.StatusUnchanged {
				continue
			}
			fmt.Fprintf(w, "  %s %s\n", r.Status, r.Key)
			printed++
		}
		if printed == 0 {
			fmt.Fprintln(w, "  (no changes)")
		}
	}
	return nil
}

func writeMatrixJSON(w io.Writer, entries []differ.MatrixEntry, _ MatrixOptions) error {
	type row struct {
		From    string         `json:"from"`
		To      string         `json:"to"`
		Results []differ.Result `json:"results"`
	}
	out := make([]row, len(entries))
	for i, e := range entries {
		out[i] = row{From: e.From, To: e.To, Results: e.Results}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
