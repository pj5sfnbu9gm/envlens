// Package counter provides utilities for counting and summarizing
// environment variable statistics across one or more env maps.
package counter

import "sort"

// Stats holds aggregated counts for an env map.
type Stats struct {
	Total     int
	Empty     int
	NonEmpty  int
	Prefixes  map[string]int
	LongestKey string
	LongestVal string
}

// Count returns Stats for the given env map.
func Count(env map[string]string) Stats {
	s := Stats{
		Prefixes: make(map[string]int),
	}
	for k, v := range env {
		s.Total++
		if v == "" {
			s.Empty++
		} else {
			s.NonEmpty++
		}
		if len(k) > len(s.LongestKey) {
			s.LongestKey = k
		}
		if len(v) > len(s.LongestVal) {
			s.LongestVal = k
		}
		if p := extractPrefix(k); p != "" {
			s.Prefixes[p]++
		}
	}
	return s
}

// TopPrefixes returns up to n prefixes sorted by count descending.
func TopPrefixes(s Stats, n int) []string {
	type kv struct {
		key   string
		count int
	}
	var pairs []kv
	for k, c := range s.Prefixes {
		pairs = append(pairs, kv{k, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	result := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, p.key)
	}
	return result
}

func extractPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' && i > 0 {
			return key[:i]
		}
	}
	return ""
}
