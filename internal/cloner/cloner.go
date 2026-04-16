// Package cloner provides functionality to deep-copy environment variable maps
// with optional key transformation applied during the clone operation.
package cloner

import "strings"

// Options configures the Clone operation.
type Options struct {
	// KeyPrefix adds a prefix to every key in the cloned map.
	KeyPrefix string
	// KeySuffix adds a suffix to every key in the cloned map.
	KeySuffix string
	// UppercaseKeys forces all keys to uppercase in the cloned map.
	UppercaseKeys bool
	// FilterKeys, when non-empty, only includes keys that start with one of these prefixes.
	FilterKeys []string
}

// DefaultOptions returns an Options with no transformations.
func DefaultOptions() Options {
	return Options{}
}

// Clone returns a deep copy of env, applying any transformations specified in opts.
func Clone(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(opts.FilterKeys) > 0 && !matchesAnyPrefix(k, opts.FilterKeys) {
			continue
		}
		newKey := k
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		newKey = opts.KeyPrefix + newKey + opts.KeySuffix
		out[newKey] = v
	}
	return out
}

func matchesAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
