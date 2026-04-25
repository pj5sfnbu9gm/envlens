package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envlens/internal/differ"
)

// ShadowReportOptions configures ReportShadow output.
type ShadowReportOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultShadowReportOptions returns sensible defaults writing to stdout.
func DefaultShadowReportOptions(w io.Writer) ShadowReportOptions {
	return ShadowReportOptions{Format: "text", Writer: w}
}

// ReportShadow renders shadow-comparison results.
func ReportShadow(entries map[string][]differ.ShadowEntry, opts ShadowReportOptions) error {
	switch opts.Format {
	case "json":
		return writeShadowJSON(entries, opts.Writer)
	default:
		return writeShadowText(entries, opts.Writer)
	}
}

func writeShadowText(entries map[string][]differ.ShadowEntry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "shadow: no discrepancies found")
		return err
	}
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		for _, e := range entries[key] {
			switch {
			case e.OnlyInShadow:
				_, err := fmt.Fprintf(w, "+ shadow-only  %s = %q\n", key, e.ShadowValue)
				if err != nil {
					return err
				}
			case e.OnlyInPrimary:
				_, err := fmt.Fprintf(w, "- primary-only %s = %q\n", key, e.PrimaryValue)
				if err != nil {
					return err
				}
			default:
				_, err := fmt.Fprintf(w, "~ changed      %s: primary=%q shadow=%q\n", key, e.PrimaryValue, e.ShadowValue)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func writeShadowJSON(entries map[string][]differ.ShadowEntry, w io.Writer) error {
	type payload struct {
		Discrepancies map[string][]differ.ShadowEntry `json:"discrepancies"`
		Total         int                             `json:"total"`
	}
	total := 0
	for _, v := range entries {
		total += len(v)
	}
	p := payload{Discrepancies: entries, Total: total}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
