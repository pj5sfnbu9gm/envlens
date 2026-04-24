package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envlens/internal/differ"
)

// DigestReportOptions configures digest report output.
type DigestReportOptions struct {
	Format  string // "text" or "json"
	Writer  io.Writer
	Verbose bool
}

// DefaultDigestReportOptions returns sensible defaults writing to stdout.
func DefaultDigestReportOptions(w io.Writer) DigestReportOptions {
	return DigestReportOptions{
		Format: "text",
		Writer: w,
	}
}

// ReportDigest writes a formatted digest report to the configured writer.
func ReportDigest(results []differ.DigestResult, opts DigestReportOptions) error {
	switch opts.Format {
	case "json":
		return writeDigestJSON(results, opts)
	default:
		return writeDigestText(results, opts)
	}
}

func writeDigestText(results []differ.DigestResult, opts DigestReportOptions) error {
	if len(results) == 0 {
		fmt.Fprintln(opts.Writer, "no targets to digest")
		return nil
	}

	conflict := differ.HasDigestConflicts(results)
	for _, r := range results {
		digest := r.Digest
		if !opts.Verbose {
			digest = digest[:12]
		}
		fmt.Fprintf(opts.Writer, "%-20s %s\n", r.Target, digest)
	}

	if conflict {
		fmt.Fprintln(opts.Writer, "\n[!] targets have diverging digests")
	} else if opts.Verbose {
		fmt.Fprintln(opts.Writer, "\n[ok] all targets share the same digest")
	}
	return nil
}

func writeDigestJSON(results []differ.DigestResult, opts DigestReportOptions) error {
	type entry struct {
		Target   string `json:"target"`
		Digest   string `json:"digest"`
		Conflict bool   `json:"conflict"`
	}

	conflict := differ.HasDigestConflicts(results)
	entries := make([]entry, len(results))
	for i, r := range results {
		entries[i] = entry{Target: r.Target, Digest: r.Digest, Conflict: conflict}
	}

	enc := json.NewEncoder(opts.Writer)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
