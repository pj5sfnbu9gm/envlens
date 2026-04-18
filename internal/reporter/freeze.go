package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envlens/internal/freezer"
)

// FreezeOptions controls freeze report output.
type FreezeOptions struct {
	Format string // "text" or "json"
	Writer io.Writer
}

// DefaultFreezeOptions returns sensible defaults.
func DefaultFreezeOptions(w io.Writer) FreezeOptions {
	return FreezeOptions{Format: "text", Writer: w}
}

// ReportFreeze writes a human- or machine-readable view of a FrozenEnv.
func ReportFreeze(f *freezer.FrozenEnv, opts FreezeOptions) error {
	switch opts.Format {
	case "json":
		return writeFreezeJSON(f, opts.Writer)
	default:
		return writeFreezeText(f, opts.Writer)
	}
}

func writeFreezeText(f *freezer.FrozenEnv, w io.Writer) error {
	keys := f.Keys()
	sort.Strings(keys)
	fmt.Fprintf(w, "Frozen environment (%d keys):\n", f.Len())
	for _, k := range keys {
		v, _ := f.Get(k)
		fmt.Fprintf(w, "  %s=%s\n", k, v)
	}
	return nil
}

func writeFreezeJSON(f *freezer.FrozenEnv, w io.Writer) error {
	type payload struct {
		Total int               `json:"total"`
		Env   map[string]string `json:"env"`
	}
	p := payload{Total: f.Len(), Env: f.ToMap()}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
