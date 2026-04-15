package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// RenameOptions configures the rename report output.
type RenameOptions struct {
	Format string // "text" or "json"
	Out    io.Writer
}

// DefaultRenameOptions returns sensible defaults.
func DefaultRenameOptions(w io.Writer) RenameOptions {
	return RenameOptions{Format: "text", Out: w}
}

// RenameEntry records a single key rename.
type RenameEntry struct {
	From string `json:"from"`
	To   string `json:"to"`
	Value string `json:"value"`
}

// ReportRename compares original and renamed maps and writes a report.
func ReportRename(original, renamed map[string]string, opts RenameOptions) error {
	var entries []RenameEntry
	for k, v := range original {
		nv, ok := renamed[k]
		if !ok {
			// Key was renamed – find new name by value match (best-effort).
			for newK, newV := range renamed {
				if newV == v {
					_, existedBefore := original[newK]
					if !existedBefore {
						entries = append(entries, RenameEntry{From: k, To: newK, Value: v})
						break
					}
				}
			}
			_ = nv
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].From < entries[j].From })

	switch opts.Format {
	case "json":
		return writeRenameJSON(opts.Out, entries)
	default:
		return writeRenameText(opts.Out, entries)
	}
}

func writeRenameText(w io.Writer, entries []RenameEntry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No keys renamed.")
		return err
	}
	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "  %s -> %s\n", e.From, e.To); err != nil {
			return err
		}
	}
	return nil
}

func writeRenameJSON(w io.Writer, entries []RenameEntry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(map[string]interface{}{"renames": entries})
}
