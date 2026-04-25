package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envlens/internal/differ"
)

// BloomReportOptions configures bloom report output.
type BloomReportOptions struct {
	Format     string // "text" or "json"
	ShowGapsOnly bool
}

// DefaultBloomReportOptions returns sensible defaults.
func DefaultBloomReportOptions() BloomReportOptions {
	return BloomReportOptions{Format: "text", ShowGapsOnly: false}
}

// ReportBloom writes a bloom presence/absence report to w.
func ReportBloom(w io.Writer, entries []differ.BloomEntry, opts BloomReportOptions) error {
	if opts.ShowGapsOnly {
		var filtered []differ.BloomEntry
		for _, e := range entries {
			if len(e.AbsentIn) > 0 {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}
	switch opts.Format {
	case "json":
		return writeBloomJSON(w, entries)
	default:
		return writeBloomText(w, entries)
	}
}

func writeBloomText(w io.Writer, entries []differ.BloomEntry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "(no bloom entries)")
		return err
	}
	for _, e := range entries {
		present := strings.Join(e.PresentIn, ", ")
		if len(e.AbsentIn) == 0 {
			_, err := fmt.Fprintf(w, "%-30s present: [%s]\n", e.Key, present)
			if err != nil {
				return err
			}
		} else {
			absent := strings.Join(e.AbsentIn, ", ")
			_, err := fmt.Fprintf(w, "%-30s present: [%s]  absent: [%s]\n", e.Key, present, absent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeBloomJSON(w io.Writer, entries []differ.BloomEntry) error {
	type row struct {
		Key       string   `json:"key"`
		PresentIn []string `json:"present_in"`
		AbsentIn  []string `json:"absent_in"`
	}
	rows := make([]row, 0, len(entries))
	for _, e := range entries {
		absent := e.AbsentIn
		if absent == nil {
			absent = []string{}
		}
		rows = append(rows, row{Key: e.Key, PresentIn: e.PresentIn, AbsentIn: absent})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
