package redactor

import "strings"

// Options controls redaction behaviour.
type Options struct {
	// Keys whose values will be fully redacted.
	Keys []string
	// Prefixes whose matching keys will be fully redacted.
	Prefixes []string
	// Placeholder replaces redacted values. Defaults to "[REDACTED]".
	Placeholder string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Prefixes:    []string{"SECRET_", "TOKEN_", "PASSWORD_", "PRIVATE_"},
		Placeholder: "[REDACTED]",
	}
}

// Result holds the original and redacted environment maps.
type Result struct {
	Original map[string]string
	Redacted map[string]string
	// RedactedKeys lists every key whose value was replaced.
	RedactedKeys []string
}

// Redact applies opts to env and returns a Result.
func Redact(env map[string]string, opts Options) Result {
	if opts.Placeholder == "" {
		opts.Placeholder = "[REDACTED]"
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	redacted := make(map[string]string, len(env))
	var redactedKeys []string

	for k, v := range env {
		if shouldRedact(k, keySet, opts.Prefixes) {
			redacted[k] = opts.Placeholder
			redactedKeys = append(redactedKeys, k)
		} else {
			redacted[k] = v
		}
	}

	return Result{
		Original:     env,
		Redacted:     redacted,
		RedactedKeys: redactedKeys,
	}
}

func shouldRedact(key string, keySet map[string]struct{}, prefixes []string) bool {
	if _, ok := keySet[key]; ok {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
