package deduper

import (
	"sort"
	"strings"
)

// Result holds the outcome of a deduplication pass.
type Result struct {
	Env     map[string]string
	Removed []string
	Kept    map[string]string // duplicate key -> kept key (case-fold mode)
}

// Options controls deduplication behaviour.
type Options struct {
	// CaseFold merges keys that differ only in case, keeping the first encountered.
	CaseFold bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{CaseFold: false}
}

// Dedupe removes duplicate keys from env according to opts.
func Dedupe(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	removed := []string{}
	kept := make(map[string]string)

	keys := sortedKeys(env)

	if opts.CaseFold {
		seen := make(map[string]string) // normalised -> original key
		for _, k := range keys {
			norm := strings.ToUpper(k)
			if orig, exists := seen[norm]; exists {
				removed = append(removed, k)
				kept[k] = orig
			} else {
				seen[norm] = k
				out[k] = env[k]
			}
		}
	} else {
		for k, v := range env {
			out[k] = v
		}
	}

	return Result{Env: out, Removed: removed, Kept: kept}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
