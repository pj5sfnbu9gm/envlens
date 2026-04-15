// Package profiler analyzes environment variable maps and produces statistical
// profiles useful for understanding the shape and composition of a configuration.
//
// A Profile includes:
//   - TotalKeys: the number of key-value pairs in the environment map.
//   - EmptyValues: keys whose values are empty or whitespace-only.
//   - SensitiveKeys: keys whose names suggest they hold sensitive data
//     (e.g. containing PASSWORD, TOKEN, SECRET, API_KEY).
//   - PrefixCounts: a frequency map of underscore-delimited key prefixes
//     (e.g. "DB" for "DB_HOST").
//   - TopPrefixes: the up-to-five most common prefixes, sorted by frequency.
//
// Usage:
//
//	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
//	p := profiler.Analyze(env)
//	fmt.Println(p.TotalKeys, p.TopPrefixes)
package profiler
