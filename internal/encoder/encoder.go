// Package encoder converts an env map into various encoded string formats
// such as base64, JSON, and CSV for embedding or transmission.
package encoder

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format specifies the output encoding format.
type Format string

const (
	FormatBase64 Format = "base64"
	FormatJSON   Format = "json"
	FormatCSV    Format = "csv"
)

// Options controls encoding behaviour.
type Options struct {
	Format    Format
	SortKeys  bool
}

// DefaultOptions returns sensible encoding defaults.
func DefaultOptions() Options {
	return Options{
		Format:   FormatJSON,
		SortKeys: true,
	}
}

// Encode serialises the env map according to opts.
func Encode(env map[string]string, opts Options) (string, error) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}

	switch opts.Format {
	case FormatJSON:
		return encodeJSON(env, keys)
	case FormatBase64:
		return encodeBase64(env, keys)
	case FormatCSV:
		return encodeCSV(env, keys)
	default:
		return "", fmt.Errorf("encoder: unknown format %q", opts.Format)
	}
}

func encodeJSON(env map[string]string, keys []string) (string, error) {
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = env[k]
	}
	b, err := json.Marshal(ordered)
	if err != nil {
		return "", fmt.Errorf("encoder: json marshal: %w", err)
	}
	return string(b), nil
}

func encodeBase64(env map[string]string, keys []string) (string, error) {
	raw, err := encodeJSON(env, keys)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(raw)), nil
}

func encodeCSV(env map[string]string, keys []string) (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)
	if err := w.Write([]string{"key", "value"}); err != nil {
		return "", fmt.Errorf("encoder: csv header: %w", err)
	}
	for _, k := range keys {
		if err := w.Write([]string{k, env[k]}); err != nil {
			return "", fmt.Errorf("encoder: csv row: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", fmt.Errorf("encoder: csv flush: %w", err)
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}
