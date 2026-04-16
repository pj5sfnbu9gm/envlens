// Package scanner inspects environment variable maps for common
// security and configuration issues.
//
// It applies a set of configurable rules to each key/value pair and
// returns structured findings with severity levels:
//
//   - "error"   — must be fixed before production deployment
//   - "warning" — likely a problem, should be reviewed
//   - "info"    — informational, may be intentional
//
// Usage:
//
//	findings := scanner.Scan(env, scanner.Options{})
//	for _, f := range findings {
//		fmt.Printf("[%s] %s\n", f.Severity, f.Message)
//	}
package scanner
