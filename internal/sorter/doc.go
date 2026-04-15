// Package sorter provides utilities for ordering environment variable keys.
//
// Keys can be sorted in ascending or descending lexicographic order. When
// GroupByPrefix is enabled, keys that share a common prefix (the portion of
// the key before the first underscore, e.g. "DB" in "DB_HOST") are grouped
// together before the secondary sort is applied within each group.
//
// Example:
//
//	env := map[string]string{"DB_HOST": "localhost", "APP_NAME": "envlens"}
//	opts := sorter.Options{Order: sorter.Ascending, GroupByPrefix: true}
//	_, keys := sorter.Sort(env, opts)
//	// keys → ["APP_NAME", "DB_HOST"]
package sorter
