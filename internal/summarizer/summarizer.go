// Package summarizer produces a human-readable summary of an env map,
// including key counts, empty values, sensitive key estimates, and top prefixes.
package summarizer

import (
	"sort"
	"strings"
)

// Summary holds aggregated statistics about an environment map.
type Summary struct {
	TotalKeys     int
	EmptyValues   int
	SensitiveKeys int
	UniqueValues  int
	TopPrefixes   []PrefixCount
}

// PrefixCount pairs a prefix with how many keys share it.
type PrefixCount struct {
	Prefix string
	Count  int
}

var sensitivePatterns = []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL"}

// Summarize analyzes env and returns a Summary.
func Summarize(env map[string]string, topN int) Summary {
	prefixCounts := map[string]int{}
	valueSet := map[string]struct{}{}
	var sensitive, empty int

	for k, v := range env {
		if v == "" {
			empty++
		} else {
			valueSet[v] = struct{}{}
		}
		if isSensitive(k) {
			sensitive++
		}
		if p := extractPrefix(k); p != "" {
			prefixCounts[p]++
		}
	}

	return Summary{
		TotalKeys:     len(env),
		EmptyValues:   empty,
		SensitiveKeys: sensitive,
		UniqueValues:  len(valueSet),
		TopPrefixes:   topPrefixes(prefixCounts, topN),
	}
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

func extractPrefix(key string) string {
	if i := strings.Index(key, "_"); i > 0 {
		return key[:i]
	}
	return ""
}

func topPrefixes(counts map[string]int, n int) []PrefixCount {
	type kv struct {
		Prefix string
		Count  int
	}
	var pairs []kv
	for p, c := range counts {
		pairs = append(pairs, kv{p, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Prefix < pairs[j].Prefix
	})
	if n > 0 && len(pairs) > n {
		pairs = pairs[:n]
	}
	out := make([]PrefixCount, len(pairs))
	for i, p := range pairs {
		out[i] = PrefixCount{p.Prefix, p.Count}
	}
	return out
}
