package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourorg/envlens/internal/differ"
)

// DefaultFingerprintReportOptions returns sensible defaults for fingerprint reporting.
func DefaultFingerprintReportOptions() FingerprintReportOptions {
	return FingerprintReportOptions{
		Format:       "text",
		ShowKeyCount: true,
	}
}

// FingerprintReportOptions controls fingerprint report output.
type FingerprintReportOptions struct {
	Format       string // "text" or "json"
	ShowKeyCount bool
}

// ReportFingerprint writes fingerprint entries to w using the given options.
func ReportFingerprint(w io.Writer, entries []differ.FingerprintEntry, opts FingerprintReportOptions) error {
	switch opts.Format {
	case "json":
		return writeFingerprintJSON(w, entries)
	default:
		return writeFingerprintText(w, entries, opts)
	}
}

func writeFingerprintText(w io.Writer, entries []differ.FingerprintEntry, opts FingerprintReportOptions) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "no fingerprint entries")
		return err
	}
	for _, e := range entries {
		line := fmt.Sprintf("%-20s  %s", e.Target, e.Fingerprint)
		if opts.ShowKeyCount {
			line += fmt.Sprintf("  (%d keys)", e.KeyCount)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeFingerprintJSON(w io.Writer, entries []differ.FingerprintEntry) error {
	type jsonEntry struct {
		Target      string `json:"target"`
		Fingerprint string `json:"fingerprint"`
		KeyCount    int    `json:"key_count"`
	}
	out := make([]jsonEntry, len(entries))
	for i, e := range entries {
		out[i] = jsonEntry{
			Target:      e.Target,
			Fingerprint: e.Fingerprint,
			KeyCount:    e.KeyCount,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
