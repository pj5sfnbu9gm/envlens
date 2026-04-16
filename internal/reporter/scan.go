package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envlens/internal/scanner"
)

// ScanOptions controls scan report output.
type ScanOptions struct {
	Format string // "text" or "json"
	Out    io.Writer
}

// DefaultScanOptions returns sensible defaults.
func DefaultScanOptions() ScanOptions {
	return ScanOptions{Format: "text", Out: os.Stdout}
}

// ReportScan writes scan findings using the given options.
func ReportScan(findings []scanner.Finding, opts ScanOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Severity != findings[j].Severity {
			return findings[i].Severity < findings[j].Severity
		}
		return findings[i].Key < findings[j].Key
	})
	switch opts.Format {
	case "json":
		return writeScanJSON(findings, opts.Out)
	default:
		return writeScanText(findings, opts.Out)
	}
}

func writeScanText(findings []scanner.Finding, w io.Writer) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "scan: no issues found")
		return err
	}
	for _, f := range findings {
		_, err := fmt.Fprintf(w, "[%s] %s: %s\n", f.Severity, f.Key, f.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeScanJSON(findings []scanner.Finding, w io.Writer) error {
	type jsonFinding struct {
		Key      string `json:"key"`
		Severity string `json:"severity"`
		Message  string `json:"message"`
	}
	out := make([]jsonFinding, len(findings))
	for i, f := range findings {
		out[i] = jsonFinding{Key: f.Key, Severity: f.Severity, Message: f.Message}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
