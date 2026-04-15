package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourorg/envlens/internal/auditor"
)

// AuditOptions controls output behaviour for audit reports.
type AuditOptions struct {
	Format string // "text" or "json"
}

// DefaultAuditOptions returns sensible defaults for audit reporting.
func DefaultAuditOptions() AuditOptions {
	return AuditOptions{Format: "text"}
}

// ReportAudit writes audit findings to w using the specified format.
func ReportAudit(w io.Writer, findings []auditor.Finding, opts AuditOptions) error {
	switch opts.Format {
	case "json":
		return writeAuditJSON(w, findings)
	default:
		return writeAuditText(w, findings)
	}
}

func writeAuditText(w io.Writer, findings []auditor.Finding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintln(w, "audit: no issues found")
		return err
	}
	for _, f := range findings {
		line := fmt.Sprintf("[%s] %s: %s\n", f.Rule, f.Key, f.Message)
		if _, err := fmt.Fprint(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeAuditJSON(w io.Writer, findings []auditor.Finding) error {
	type jsonFinding struct {
		Key     string `json:"key"`
		Rule    string `json:"rule"`
		Message string `json:"message"`
	}

	out := make([]jsonFinding, 0, len(findings))
	for _, f := range findings {
		out = append(out, jsonFinding{
			Key:     f.Key,
			Rule:    f.Rule,
			Message: f.Message,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
