package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ExportOptions controls how the export report is rendered.
type ExportOptions struct {
	Format string // "text" or "json"
	Target string // target label shown in the report
}

// DefaultExportOptions returns sensible defaults for export reporting.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{Format: "text", Target: "env"}
}

// ReportExport writes a human-readable or JSON summary of exported env vars.
func ReportExport(w io.Writer, env map[string]string, opts ExportOptions) error {
	switch opts.Format {
	case "json":
		return writeExportJSON(w, env, opts)
	default:
		return writeExportText(w, env, opts)
	}
}

func writeExportText(w io.Writer, env map[string]string, opts ExportOptions) error {
	keys := sortedKeys(env)
	fmt.Fprintf(w, "Exported %d variable(s) for target: %s\n", len(keys), opts.Target)
	for _, k := range keys {
		fmt.Fprintf(w, "  %s=%s\n", k, env[k])
	}
	return nil
}

func writeExportJSON(w io.Writer, env map[string]string, opts ExportOptions) error {
	type payload struct {
		Target    string            `json:"target"`
		Count     int               `json:"count"`
		Variables map[string]string `json:"variables"`
	}
	p := payload{
		Target:    opts.Target,
		Count:     len(env),
		Variables: env,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
