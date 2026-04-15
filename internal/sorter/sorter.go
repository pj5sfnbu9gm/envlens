package sorter

import (
	"sort"
	"strings"
)

// Order defines the sort order for environment variable keys.
type Order int

const (
	// Ascending sorts keys A→Z.
	Ascending Order = iota
	// Descending sorts keys Z→A.
	Descending
)

// Options controls how the sorter behaves.
type Options struct {
	Order      Order
	GroupByPrefix bool // group keys sharing the same prefix together
}

// DefaultOptions returns sensible default sort options.
func DefaultOptions() Options {
	return Options{
		Order:         Ascending,
		GroupByPrefix: false,
	}
}

// Sort returns a new map with the same key/value pairs and a sorted slice of
// keys according to the provided options.
func Sort(env map[string]string, opts Options) (map[string]string, []string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	if opts.GroupByPrefix {
		sort.Slice(keys, func(i, j int) bool {
			pi := extractPrefix(keys[i])
			pj := extractPrefix(keys[j])
			if pi != pj {
				if opts.Order == Descending {
					return pi > pj
				}
				return pi < pj
			}
			if opts.Order == Descending {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		})
	} else {
		sort.Slice(keys, func(i, j int) bool {
			if opts.Order == Descending {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		})
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	return out, keys
}

// extractPrefix returns the portion of a key before the first underscore, or
// the entire key if no underscore is present.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
