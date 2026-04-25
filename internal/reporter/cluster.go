package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// ClusterReportOptions controls the cluster report output.
type ClusterReportOptions struct {
	Format  string // "text" or "json"
	MaskVal bool   // mask the shared value in output
}

// DefaultClusterReportOptions returns sensible defaults.
func DefaultClusterReportOptions() ClusterReportOptions {
	return ClusterReportOptions{Format: "text"}
}

// ReportCluster writes cluster entries to w.
func ReportCluster(w io.Writer, entries []differ.ClusterEntry, opts ClusterReportOptions) error {
	switch strings.ToLower(opts.Format) {
	case "json":
		return writeClusterJSON(w, entries, opts)
	default:
		return writeClusterText(w, entries, opts)
	}
}

func writeClusterText(w io.Writer, entries []differ.ClusterEntry, opts ClusterReportOptions) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "no clusters found")
		return err
	}
	for _, e := range entries {
		val := e.Value
		if opts.MaskVal {
			val = "***"
		}
		_, err := fmt.Fprintf(w, "value=%q  keys=[%s]  targets=[%s]\n",
			val,
			strings.Join(e.Keys, ", "),
			strings.Join(e.Targets, ", "),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeClusterJSON(w io.Writer, entries []differ.ClusterEntry, opts ClusterReportOptions) error {
	type row struct {
		Value   string   `json:"value"`
		Keys    []string `json:"keys"`
		Targets []string `json:"targets"`
	}
	rows := make([]row, len(entries))
	for i, e := range entries {
		v := e.Value
		if opts.MaskVal {
			v = "***"
		}
		rows[i] = row{Value: v, Keys: e.Keys, Targets: e.Targets}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
