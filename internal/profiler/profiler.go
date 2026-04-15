// Package profiler provides functionality for profiling environment variable
// sets by computing statistics such as key counts, empty value ratios,
// sensitive key counts, and prefix distributions.
package profiler

import (
	"sort"
	"strings"
)

// Profile holds statistical information about an environment variable map.
type Profile struct {
	TotalKeys     int            `json:"total_keys"`
	EmptyValues   int            `json:"empty_values"`
	SensitiveKeys int            `json:"sensitive_keys"`
	PrefixCounts  map[string]int `json:"prefix_counts"`
	TopPrefixes   []string       `json:"top_prefixes"`
}

// sensitivePatterns are substrings that indicate a key may be sensitive.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL",
}

// Analyze computes a Profile for the given environment map.
func Analyze(env map[string]string) Profile {
	p := Profile{
		PrefixCounts: make(map[string]int),
	}

	for k, v := range env {
		p.TotalKeys++

		if strings.TrimSpace(v) == "" {
			p.EmptyValues++
		}

		if isSensitive(k) {
			p.SensitiveKeys++
		}

		if prefix := extractPrefix(k); prefix != "" {
			p.PrefixCounts[prefix]++
		}
	}

	p.TopPrefixes = topN(p.PrefixCounts, 5)
	return p
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pat := range sensitivePatterns {
		if strings.Contains(upper, pat) {
			return true
		}
	}
	return false
}

func extractPrefix(key string) string {
	parts := strings.SplitN(key, "_", 2)
	if len(parts) == 2 && parts[0] != "" {
		return parts[0]
	}
	return ""
}

func topN(counts map[string]int, n int) []string {
	type kv struct {
		Key   string
		Count int
	}
	var pairs []kv
	for k, v := range counts {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Key < pairs[j].Key
	})
	var result []string
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, p.Key)
	}
	return result
}
