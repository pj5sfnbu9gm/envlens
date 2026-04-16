package typecast

import (
	"fmt"
	"strconv"
	"strings"
)

// Result holds the outcome of casting a single env var.
type Result struct {
	Key      string
	Raw      string
	Kind     string
	Value    interface{}
	Error    string
}

// Options controls how env vars are cast.
type Options struct {
	// Hints maps key names to desired types: "int", "float", "bool", "string"
	Hints map[string]string
	// StrictMode causes Cast to return an error on first cast failure.
	StrictMode bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Hints:      map[string]string{},
		StrictMode: false,
	}
}

// Cast attempts to cast each env var according to the provided hints.
// Keys without hints are returned as strings.
func Cast(env map[string]string, opts Options) ([]Result, error) {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		kind, ok := opts.Hints[k]
		if !ok {
			kind = "string"
		}
		r := Result{Key: k, Raw: v, Kind: kind}
		switch strings.ToLower(kind) {
		case "int":
			n, err := strconv.Atoi(strings.TrimSpace(v))
			if err != nil {
				r.Error = fmt.Sprintf("cannot cast %q to int", v)
				if opts.StrictMode {
					return nil, fmt.Errorf("typecast: key %s: %s", k, r.Error)
				}
			} else {
				r.Value = n
			}
		case "float":
			f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err != nil {
				r.Error = fmt.Sprintf("cannot cast %q to float", v)
				if opts.StrictMode {
					return nil, fmt.Errorf("typecast: key %s: %s", k, r.Error)
				}
			} else {
				r.Value = f
			}
		case "bool":
			b, err := strconv.ParseBool(strings.TrimSpace(v))
			if err != nil {
				r.Error = fmt.Sprintf("cannot cast %q to bool", v)
				if opts.StrictMode {
					return nil, fmt.Errorf("typecast: key %s: %s", k, r.Error)
				}
			} else {
				r.Value = b
			}
		default:
			r.Value = v
		}
		results = append(results, r)
	}
	return results, nil
}
