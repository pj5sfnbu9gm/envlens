package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yourusername/envlens/internal/differ"
)

// CensusReportOptions controls census report output.
type CensusReportOptions struct {
	Format       string // "text" or "json"
	Out          io.Writer
	TotalTargets int
	ShowGapsOnly bool
}

// DefaultCensusReportOptions returns sensible defaults.
func DefaultCensusReportOptions(totalTargets int) CensusReportOptions {
	return CensusReportOptions{
		Format:       "text",
		Out:          os.Stdout,
		TotalTargets: totalTargets,
		ShowGapsOnly: false,
	}
}

// ReportCensus writes census entries to the configured output.
func ReportCensus(entries []differ.CensusEntry, opts CensusReportOptions) error {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	switch opts.Format {
	case "json":
		return writeCensusJSON(entries, opts)
	default:
		return writeCensusText(entries, opts)
	}
}

func writeCensusText(entries []differ.CensusEntry, opts CensusReportOptions) error {
	if len(entries) == 0 {
		fmt.Fprintln(opts.Out, "no census data")
		return nil
	}
	for _, e := range entries {
		if opts.ShowGapsOnly && e.Count == opts.TotalTargets {
			continue
		}
		gap := ""
		if e.Count < opts.TotalTargets {
			gap = " [GAP]"
		}
		fmt.Fprintf(opts.Out, "%-40s %d/%d (%.0f%%)%s\n",
			e.Key, e.Count, opts.TotalTargets, e.Coverage*100, gap)
	}
	return nil
}

func writeCensusJSON(entries []differ.CensusEntry, opts CensusReportOptions) error {
	type row struct {
		Key      string   `json:"key"`
		Count    int      `json:"count"`
		Total    int      `json:"total"`
		Coverage float64  `json:"coverage"`
		Targets  []string `json:"targets"`
		Gap      bool     `json:"gap"`
	}
	var rows []row
	for _, e := range entries {
		if opts.ShowGapsOnly && e.Count == opts.TotalTargets {
			continue
		}
		rows = append(rows, row{
			Key:      e.Key,
			Count:    e.Count,
			Total:    opts.TotalTargets,
			Coverage: e.Coverage,
			Targets:  e.Targets,
			Gap:      e.Count < opts.TotalTargets,
		})
	}
	enc := json.NewEncoder(opts.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
