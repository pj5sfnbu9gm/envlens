// Package summarizer provides statistical analysis of environment variable maps.
//
// It computes metrics such as total key count, empty value count, estimated
// sensitive key count, unique value count, and the most common key prefixes.
//
// Example:
//
//	env := map[string]string{
//		"DB_HOST":     "localhost",
//		"DB_PASSWORD": "s3cr3t",
//		"APP_ENV":     "production",
//	}
//	summary := summarizer.Summarize(env, 5)
//	fmt.Println(summary.TotalKeys)     // 3
//	fmt.Println(summary.SensitiveKeys) // 1
package summarizer
