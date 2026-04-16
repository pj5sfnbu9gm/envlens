package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// TemplateOptions controls how template render results are reported.
type TemplateOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultTemplateOptions returns sensible defaults writing to stdout.
func DefaultTemplateOptions(w io.Writer) TemplateOptions {
	return TemplateOptions{Format: "text", Writer: w}
}

// TemplateResult bundles the rendered output with metadata for reporting.
type TemplateResult struct {
	Output      string
	MissingKeys []string
}

// ReportTemplate writes a template render result to the configured writer.
func ReportTemplate(res TemplateResult, opts TemplateOptions) error {
	switch opts.Format {
	case "json":
		return writeTemplateJSON(opts.Writer, res)
	default:
		return writeTemplateText(opts.Writer, res)
	}
}

func writeTemplateText(w io.Writer, res TemplateResult) error {
	fmt.Fprintln(w, res.Output)
	if len(res.MissingKeys) > 0 {
		sorted := append([]string(nil), res.MissingKeys...)
		sort.Strings(sorted)
		fmt.Fprintln(w, "\n[missing keys]")
		for _, k := range sorted {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}
	return nil
}

func writeTemplateJSON(w io.Writer, res TemplateResult) error {
	payload := map[string]interface{}{
		"output":       res.Output,
		"missing_keys": res.MissingKeys,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
