// Package linter provides a rule-based linting engine for environment variable
// maps. It evaluates each key-value pair against a configurable set of Rule
// functions and collects all findings with associated severity levels.
//
// Built-in rules cover common issues such as empty values, lowercase keys,
// keys containing whitespace, and keys prefixed with an underscore. Custom
// rules can be supplied to Lint alongside or in place of DefaultRules.
//
// Example usage:
//
//	env := map[string]string{"db_host": "", "APP_PORT": "8080"}
//	findings := linter.Lint(env, linter.DefaultRules())
//	for _, f := range findings {
//		fmt.Printf("[%s] %s: %s\n", f.Severity, f.Key, f.Message)
//	}
package linter
