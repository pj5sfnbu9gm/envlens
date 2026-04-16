package templater

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Options controls template rendering behavior.
type Options struct {
	// FailOnMissing causes Render to return an error if a key is referenced but not in env.
	FailOnMissing bool
	// LeftDelim and RightDelim override the default {{ }} delimiters.
	LeftDelim  string
	RightDelim string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: false,
		LeftDelim:     "{{",
		RightDelim:    "}}",
	}
}

// Result holds the rendered output alongside metadata.
type Result struct {
	Output      string
	MissingKeys []string
}

// Render applies env values to the given template text.
func Render(tmplText string, env map[string]string, opts Options) (Result, error) {
	var missing []string

	funcMap := template.FuncMap{
		"env": func(key string) (string, error) {
			v, ok := env[key]
			if !ok {
				missing = append(missing, key)
				if opts.FailOnMissing {
					return "", fmt.Errorf("missing env key: %s", key)
				}
				return "", nil
			}
			return v, nil
		},
		"envOr": func(key, fallback string) string {
			if v, ok := env[key]; ok {
				return v
			}
			return fallback
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	t, err := template.New("envlens").
		Delims(opts.LeftDelim, opts.RightDelim).
		Funcs(funcMap).
		Parse(tmplText)
	if err != nil {
		return Result{}, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, env); err != nil {
		return Result{}, fmt.Errorf("execute template: %w", err)
	}

	return Result{Output: buf.String(), MissingKeys: missing}, nil
}
