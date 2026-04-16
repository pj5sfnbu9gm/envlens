package scorer

import (
	"sort"
	"strings"
)

// Score holds the result of evaluating an environment map.
type Score struct {
	MaxScore   int
	FinalScore int
	Penalties  []Penalty
}

// Penalty is a single deduction applied during scoring.
type Penalty struct {
	Key    string
	Reason string
	Points int
}

// Options controls scoring behaviour.
type Options struct {
	EmptyValuePenalty   int
	LowercaseKeyPenalty int
	WhitespacePenalty   int
	MaxScore            int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		EmptyValuePenalty:   5,
		LowercaseKeyPenalty: 3,
		WhitespacePenalty:   4,
		MaxScore:            100,
	}
}

// ScoreEnv evaluates the quality of an env map and returns a Score.
func ScoreEnv(env map[string]string, opts Options) Score {
	var penalties []Penalty

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		if v == "" {
			penalties = append(penalties, Penalty{Key: k, Reason: "empty value", Points: opts.EmptyValuePenalty})
		}
		if k != strings.ToUpper(k) {
			penalties = append(penalties, Penalty{Key: k, Reason: "lowercase key", Points: opts.LowercaseKeyPenalty})
		}
		if strings.ContainsAny(k, " \t") {
			penalties = append(penalties, Penalty{Key: k, Reason: "whitespace in key", Points: opts.WhitespacePenalty})
		}
	}

	final := opts.MaxScore
	for _, p := range penalties {
		final -= p.Points
	}
	if final < 0 {
		final = 0
	}

	return Score{
		MaxScore:   opts.MaxScore,
		FinalScore: final,
		Penalties:  penalties,
	}
}
