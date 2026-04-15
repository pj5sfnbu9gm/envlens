// Package auditor provides rule-based auditing of environment variable maps.
//
// It supports built-in rules and custom user-defined rules. Each rule
// receives a key-value pair and returns an error if the pair violates
// the rule's constraints.
//
// Built-in rules include:
//
//   - no-empty-value: flags keys with empty string values
//   - no-whitespace-key: flags keys containing spaces or tabs
//   - uppercase-key: flags keys that contain lowercase letters
//
// Example usage:
//
//	env := map[string]string{"APP_ENV": "production", "db_url": ""}
//	findings := auditor.Audit(env, auditor.DefaultRules())
//	for _, f := range findings {
//		fmt.Printf("[%s] %s: %s\n", f.Rule, f.Key, f.Message)
//	}
package auditor
