package normalizer

import (
	"strings"
)

// Options controls normalization behavior.
type Options struct {
	UppercaseKeys   bool
	TrimSpace       bool
	RemoveEmpty     bool
	ReplaceHyphens  bool // replace - with _ in keys
}

// DefaultOptions returns sensible normalization defaults.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:  true,
		TrimSpace:      true,
		RemoveEmpty:    false,
		ReplaceHyphens: true,
	}
}

// Result holds the original and normalized key/value plus a flag if changed.
type Result struct {
	OriginalKey   string
	NormalizedKey string
	OriginalVal   string
	NormalizedVal string
	Changed       bool
}

// Normalize applies normalization rules to an env map and returns results.
func Normalize(env map[string]string, opts Options) (map[string]string, []Result) {
	out := make(map[string]string, len(env))
	results := make([]Result, 0, len(env))

	for k, v := range env {
		nk := k
		nv := v

		if opts.ReplaceHyphens {
			nk = strings.ReplaceAll(nk, "-", "_")
		}
		if opts.UppercaseKeys {
			nk = strings.ToUpper(nk)
		}
		if opts.TrimSpace {
			nk = strings.TrimSpace(nk)
			nv = strings.TrimSpace(nv)
		}

		changed := nk != k || nv != v

		if opts.RemoveEmpty && nv == "" {
			continue
		}

		out[nk] = nv
		results = append(results, Result{
			OriginalKey:   k,
			NormalizedKey: nk,
			OriginalVal:   v,
			NormalizedVal: nv,
			Changed:       changed,
		})
	}

	return out, results
}
