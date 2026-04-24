package differ

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// DigestResult holds the computed digest for a single target.
type DigestResult struct {
	Target string
	Digest string
}

// DigestOptions controls how digests are computed.
type DigestOptions struct {
	// IncludeKeys limits digesting to specific keys. If empty, all keys are included.
	IncludeKeys []string
}

// DefaultDigestOptions returns sensible defaults.
func DefaultDigestOptions() DigestOptions {
	return DigestOptions{}
}

// Digest computes a deterministic SHA-256 fingerprint for each target's env map.
// Keys are sorted before hashing to ensure consistency regardless of insertion order.
func Digest(targets map[string]map[string]string, opts DigestOptions) []DigestResult {
	results := make([]DigestResult, 0, len(targets))

	names := make([]string, 0, len(targets))
	for name := range targets {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		env := targets[name]
		keys := keysForDigest(env, opts.IncludeKeys)
		sort.Strings(keys)

		h := sha256.New()
		for _, k := range keys {
			fmt.Fprintf(h, "%s=%s\n", k, env[k])
		}

		results = append(results, DigestResult{
			Target: name,
			Digest: hex.EncodeToString(h.Sum(nil)),
		})
	}

	return results
}

// HasDigestConflicts returns true if any two targets produce different digests.
func HasDigestConflicts(results []DigestResult) bool {
	if len(results) < 2 {
		return false
	}
	first := results[0].Digest
	for _, r := range results[1:] {
		if r.Digest != first {
			return true
		}
	}
	return false
}

func keysForDigest(env map[string]string, include []string) []string {
	if len(include) == 0 {
		keys := make([]string, 0, len(env))
		for k := range env {
			keys = append(keys, k)
		}
		return keys
	}
	keys := make([]string, 0, len(include))
	for _, k := range include {
		if _, ok := env[k]; ok {
			keys = append(keys, k)
		}
	}
	return keys
}
