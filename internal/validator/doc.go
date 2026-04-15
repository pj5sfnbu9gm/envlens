// Package validator provides rule-based validation for environment variable
// maps loaded by envlens.
//
// It ships with a set of DefaultRules covering common pitfalls such as empty
// values, improperly formatted keys, and accidental leading/trailing
// whitespace. Callers may supply additional custom rules to enforce
// project-specific conventions.
//
// Usage:
//
//	env, _ := loader.LoadFile("production.env")
//	findings := validator.Validate(env, validator.DefaultRules())
//	for _, f := range findings {
//		fmt.Println(f.Message)
//	}
package validator
