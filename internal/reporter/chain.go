package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/user/envlens/internal/differ"
)

// ChainOptions controls chain diff report output.
type ChainOptions struct {
	Writer io.Writer
	Format string // "text" or "json"
	ShowUnchanged bool
}

// DefaultChainOptions returns sensible defaults.
func DefaultChainOptions() ChainOptions {
	return ChainOptions{Writer: os.Stdout, Format: "text"}
}

// ReportChain writes a chain diff report to the configured writer.
func ReportChain(chain []differ.ChainResult, opts ChainOptions) error {
	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}
	switch opts.Format {
	case "json":
		return writeChainJSON(chain, opts)
	default:
		return writeChainText(chain, opts)
	}
}

func writeChainText(chain []differ.ChainResult, opts ChainOptions) error {
	if len(chain) == 0 {
		fmt.Fprintln(opts.Writer, "no chain steps to report")
		return nil
	}
	for _, step := range chain {
		fmt.Fprintf(opts.Writer, "--- %s -> %s ---\n", step.From, step.To)
		for _, r := range step.Results {
			if r.Status == "unchanged" && !opts.ShowUnchanged {
				continue
			}
			fmt.Fprintf(opts.Writer, "  [%s] %s\n", r.Status, r.Key)
		}
	}
	return nil
}

func writeChainJSON(chain []differ.ChainResult, opts ChainOptions) error {
	enc := json.NewEncoder(opts.Writer)
	enc.SetIndent("", "  ")
	return enc.Encode(chain)
}
