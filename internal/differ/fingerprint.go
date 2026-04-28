package differ

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// FingerprintEntry holds the computed fingerprint for a single target.
type FingerprintEntry struct {
	Target      string
	Fingerprint string
	KeyCount    int
}

// DefaultFingerprintOptions returns a FingerprintOptions with sensible defaults.
func DefaultFingerprintOptions() FingerprintOptions {
	return FingerprintOptions{
		IncludeValues: true,
	}
}

// FingerprintOptions controls how fingerprints are computed.
type FingerprintOptions struct {
	// IncludeValues includes env values in the hash when true.
	// When false, only keys are hashed (useful for structural comparison).
	IncludeValues bool
	// OnlyKeys restricts hashing to the provided key names.
	OnlyKeys []string
}

// Fingerprint computes a deterministic SHA-256 hash for each target's
// environment map and returns a slice of FingerprintEntry.
func Fingerprint(targets map[string]map[string]string, opts FingerprintOptions) []FingerprintEntry {
	if len(targets) == 0 {
		return nil
	}

	filterSet := make(map[string]struct{}, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		filterSet[k] = struct{}{}
	}

	names := make([]string, 0, len(targets))
	for name := range targets {
		names = append(names, name)
	}
	sort.Strings(names)

	results := make([]FingerprintEntry, 0, len(names))
	for _, name := range names {
		env := targets[name]
		keys := make([]string, 0, len(env))
		for k := range env {
			if len(filterSet) > 0 {
				if _, ok := filterSet[k]; !ok {
					continue
				}
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var sb strings.Builder
		for _, k := range keys {
			sb.WriteString(k)
			if opts.IncludeValues {
				sb.WriteByte('=')
				sb.WriteString(env[k])
			}
			sb.WriteByte('\n')
		}

		sum := sha256.Sum256([]byte(sb.String()))
		results = append(results, FingerprintEntry{
			Target:      name,
			Fingerprint: fmt.Sprintf("%x", sum),
			KeyCount:    len(keys),
		})
	}
	return results
}

// HasFingerprintConflicts returns true when any two entries share the same
// fingerprint but have different target names — indicating an exact duplicate.
func HasFingerprintConflicts(entries []FingerprintEntry) bool {
	seen := make(map[string]string, len(entries))
	for _, e := range entries {
		if prev, ok := seen[e.Fingerprint]; ok && prev != e.Target {
			return true
		}
		seen[e.Fingerprint] = e.Target
	}
	return false
}
