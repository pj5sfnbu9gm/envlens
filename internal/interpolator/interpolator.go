package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR} and $VAR style references.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// FallbackToOS allows unresolved references to be looked up in os.Environ.
	FallbackToOS bool
	// FailOnMissing returns an error when a referenced variable cannot be resolved.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FallbackToOS:  true,
		FailOnMissing: false,
	}
}

// Interpolate expands variable references inside values of env using the
// variables already present in env (and optionally os.Environ).
// It returns a new map with all values expanded.
func Interpolate(env map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		expanded, err := expand(v, env, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolator: key %q: %w", k, err)
		}
		result[k] = expanded
	}
	return result, nil
}

func expand(value string, env map[string]string, opts Options) (string, error) {
	var expandErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := extractName(match)
		if v, ok := env[name]; ok {
			return v
		}
		if opts.FallbackToOS {
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
		}
		if opts.FailOnMissing {
			expandErr = fmt.Errorf("unresolved variable %q", name)
			return match
		}
		return ""
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
