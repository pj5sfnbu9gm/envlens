package renamer

import (
	"fmt"
	"strings"
)

// Rule defines a single rename transformation.
type Rule struct {
	From string
	To   string
}

// Options controls renaming behaviour.
type Options struct {
	// Rules is the ordered list of rename rules to apply.
	Rules []Rule
	// FailOnMissing causes Rename to return an error if a From key is not found.
	FailOnMissing bool
	// Prefix replaces a key prefix when OldPrefix is non-empty.
	OldPrefix string
	NewPrefix string
}

// DefaultOptions returns a zero-value Options ready for use.
func DefaultOptions() Options {
	return Options{}
}

// Rename applies renaming rules to env and returns a new map.
// The original map is never mutated.
func Rename(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	// Apply explicit key rules first.
	for _, rule := range opts.Rules {
		if rule.From == "" || rule.To == "" {
			return nil, fmt.Errorf("renamer: rule has empty From or To field")
		}
		val, ok := out[rule.From]
		if !ok {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("renamer: key %q not found in env", rule.From)
			}
			continue
		}
		delete(out, rule.From)
		out[rule.To] = val
	}

	// Apply prefix replacement.
	if opts.OldPrefix != "" {
		for k, v := range out {
			if strings.HasPrefix(k, opts.OldPrefix) {
				newKey := opts.NewPrefix + strings.TrimPrefix(k, opts.OldPrefix)
				delete(out, k)
				out[newKey] = v
			}
		}
	}

	return out, nil
}
