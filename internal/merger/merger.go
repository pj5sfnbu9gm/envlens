// Package merger provides utilities for merging multiple environment
// variable maps into a single unified map, with configurable precedence
// and conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how conflicts are resolved when the same key exists
// in multiple sources.
type Strategy int

const (
	// StrategyFirst keeps the value from the first source that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last source that defines the key.
	StrategyLast
	// StrategyError returns an error if the same key appears in multiple sources.
	StrategyError
)

// Conflict records a key that appeared in more than one source.
type Conflict struct {
	Key    string
	First  string
	Second string
}

// Result holds the merged environment map and any conflicts encountered.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Merge combines the provided named environment maps according to the given
// strategy. Sources are processed in the order they are supplied.
func Merge(sources map[string]map[string]string, order []string, strategy Strategy) (*Result, error) {
	result := &Result{
		Env:       make(map[string]string),
		Conflicts: []Conflict{},
	}

	origin := make(map[string]string) // tracks which source set each key

	for _, name := range order {
		env, ok := sources[name]
		if !ok {
			continue
		}
		for k, v := range env {
			if existing, seen := result.Env[k]; seen {
				conflict := Conflict{Key: k, First: existing, Second: v}
				switch strategy {
				case StrategyError:
					return nil, fmt.Errorf("merger: key %q defined in both %q and %q", k, origin[k], name)
				case StrategyLast:
					result.Env[k] = v
					origin[k] = name
					result.Conflicts = append(result.Conflicts, conflict)
				case StrategyFirst:
					result.Conflicts = append(result.Conflicts, conflict)
				}
			} else {
				result.Env[k] = v
				origin[k] = name
			}
		}
	}

	return result, nil
}
