package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/envlens/internal/typecast"
)

// TypecastOptions controls output of ReportTypecast.
type TypecastOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultTypecastOptions returns sensible defaults.
func DefaultTypecastOptions(w io.Writer) TypecastOptions {
	return TypecastOptions{Format: "text", Writer: w}
}

// ReportTypecast writes a formatted report of typecast results.
func ReportTypecast(results []typecast.Result, opts TypecastOptions) error {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	switch opts.Format {
	case "json":
		return writeTypecastJSON(results, opts.Writer)
	default:
		return writeTypecastText(results, opts.Writer)
	}
}

func writeTypecastText(results []typecast.Result, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tTYPE\tVALUE\tERROR")
	for _, r := range results {
		val := fmt.Sprintf("%v", r.Value)
		if r.Error != "" {
			val = ""
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Key, r.Kind, val, r.Error)
	}
	return tw.Flush()
}

func writeTypecastJSON(results []typecast.Result, w io.Writer) error {
	type row struct {
		Key   string      `json:"key"`
		Kind  string      `json:"kind"`
		Value interface{} `json:"value,omitempty"`
		Error string      `json:"error,omitempty"`
	}
	rows := make([]row, len(results))
	for i, r := range results {
		rows[i] = row{Key: r.Key, Kind: r.Kind, Value: r.Value, Error: r.Error}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
