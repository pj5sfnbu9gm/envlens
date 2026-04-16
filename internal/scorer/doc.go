// Package scorer provides quality scoring for environment variable maps.
//
// ScoreEnv evaluates a map[string]string against a set of rules and returns
// a Score containing a final numeric grade (0–MaxScore) and a list of
// Penalty entries that explain each deduction.
//
// Example:
//
//	result := scorer.ScoreEnv(env, scorer.DefaultOptions())
//	fmt.Printf("Score: %d/%d\n", result.FinalScore, result.MaxScore)
package scorer
