package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/envlens/internal/masker"
)

// MaskedOptions configures the masked-env report.
type MaskedOptions struct {
	Format string // "text" or "json"
	Mask   masker.MaskOptions
}

// DefaultMaskedOptions returns sensible defaults for a masked report.
func DefaultMaskedOptions() MaskedOptions {
	return MaskedOptions{
		Format: "text",
		Mask:   masker.DefaultMaskOptions(),
	}
}

// ReportMasked writes a masked view of an environment map to w.
func ReportMasked(w io.Writer, env map[string]string, opts MaskedOptions) error {
	masked := masker.MaskEnv(env, opts.Mask)

	switch opts.Format {
	case "json":
		return writeMaskedJSON(w, masked)
	default:
		return writeMaskedText(w, masked)
	}
}

func writeMaskedText(w io.Writer, env map[string]string) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, env[k]); err != nil {
			return err
		}
	}
	return nil
}

func writeMaskedJSON(w io.Writer, env map[string]string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(env)
}
