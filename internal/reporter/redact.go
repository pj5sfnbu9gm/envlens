package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envlens/internal/redactor"
)

// RedactOptions controls how redaction reports are rendered.
type RedactOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultRedactOptions returns text output to stdout.
func DefaultRedactOptions(w io.Writer) RedactOptions {
	return RedactOptions{Format: "text", Writer: w}
}

// ReportRedact writes a human- or machine-readable summary of a redaction result.
func ReportRedact(result redactor.Result, opts RedactOptions) error {
	switch opts.Format {
	case "json":
		return writeRedactJSON(result, opts.Writer)
	default:
		return writeRedactText(result, opts.Writer)
	}
}

func writeRedactText(result redactor.Result, w io.Writer) error {
	keys := make([]string, 0, len(result.Redacted))
	for k := range result.Redacted {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	redactedSet := make(map[string]struct{}, len(result.RedactedKeys))
	for _, k := range result.RedactedKeys {
		redactedSet[k] = struct{}{}
	}

	for _, k := range keys {
		marker := " "
		if _, ok := redactedSet[k]; ok {
			marker = "*"
		}
		fmt.Fprintf(w, "[%s] %s=%s\n", marker, k, result.Redacted[k])
	}
	if len(result.RedactedKeys) == 0 {
		fmt.Fprintln(w, "No keys redacted.")
	}
	return nil
}

func writeRedactJSON(result redactor.Result, w io.Writer) error {
	type payload struct {
		Redacted     map[string]string `json:"redacted"`
		RedactedKeys []string          `json:"redacted_keys"`
	}
	p := payload{
		Redacted:     result.Redacted,
		RedactedKeys: result.RedactedKeys,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
