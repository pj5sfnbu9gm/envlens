package exporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format represents the export output format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatExport Format = "export"
)

// Options configures the export behaviour.
type Options struct {
	Format  Format
	Sorted  bool
	Comment string // optional header comment
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		Format: FormatDotenv,
		Sorted: true,
	}
}

// Export writes env vars from env to w using the given options.
func Export(w io.Writer, env map[string]string, opts Options) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.Sorted {
		sort.Strings(keys)
	}

	if opts.Comment != "" {
		fmt.Fprintf(w, "# %s\n", opts.Comment)
	}

	for _, k := range keys {
		v := env[k]
		switch opts.Format {
		case FormatShell:
			fmt.Fprintf(w, "%s=%q\n", k, v)
		case FormatExport:
			fmt.Fprintf(w, "export %s=%q\n", k, v)
		default: // dotenv
			if needsQuotes(v) {
				fmt.Fprintf(w, "%s=%q\n", k, v)
			} else {
				fmt.Fprintf(w, "%s=%s\n", k, v)
			}
		}
	}
	return nil
}

// ExportToFile writes the env map to the given file path.
func ExportToFile(path string, env map[string]string, opts Options) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("exporter: create file %q: %w", path, err)
	}
	defer f.Close()
	return Export(f, env, opts)
}

func needsQuotes(v string) bool {
	return strings.ContainsAny(v, " \t\n#$'\"\\")
}
