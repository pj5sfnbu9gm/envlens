package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// WindowReportOptions configures output for ReportWindow.
type WindowReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultWindowReportOptions returns text-format defaults writing to stdout.
func DefaultWindowReportOptions(w io.Writer) WindowReportOptions {
	return WindowReportOptions{Format: "text", Writer: w}
}

// ReportWindow renders sliding-window diff results.
func ReportWindow(windows []differ.WindowResult, opts WindowReportOptions) error {
	switch strings.ToLower(opts.Format) {
	case "json":
		return writeWindowJSON(windows, opts.Writer)
	default:
		return writeWindowText(windows, opts.Writer)
	}
}

func writeWindowText(windows []differ.WindowResult, w io.Writer) error {
	if len(windows) == 0 {
		_, err := fmt.Fprintln(w, "no windows to display")
		return err
	}
	for _, win := range windows {
		label := strings.Join(win.Labels, " → ")
		fmt.Fprintf(w, "[window] %s\n", label)
		if len(win.Results) == 0 {
			fmt.Fprintln(w, "  (no changes)")
			continue
		}
		for _, r := range win.Results {
			switch r.Status {
			case "added":
				fmt.Fprintf(w, "  + %s=%s\n", r.Key, r.NewValue)
			case "removed":
				fmt.Fprintf(w, "  - %s=%s\n", r.Key, r.OldValue)
			case "changed":
				fmt.Fprintf(w, "  ~ %s: %s → %s\n", r.Key, r.OldValue, r.NewValue)
			default:
				fmt.Fprintf(w, "    %s=%s\n", r.Key, r.NewValue)
			}
		}
	}
	return nil
}

func writeWindowJSON(windows []differ.WindowResult, w io.Writer) error {
	type entry struct {
		Labels  []string        `json:"labels"`
		Results []differ.Result `json:"results"`
	}
	out := make([]entry, 0, len(windows))
	for _, win := range windows {
		out = append(out, entry{Labels: win.Labels, Results: win.Results})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
