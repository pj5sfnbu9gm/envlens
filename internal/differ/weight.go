package differ

// Weight assigns a numeric importance score to each diff result based on
// the kind of change and optional per-key multipliers.

// WeightOptions controls how scores are computed.
type WeightOptions struct {
	// Scores per status kind.
	AddedScore   float64
	RemovedScore float64
	ChangedScore float64

	// KeyWeights overrides the base score for specific keys.
	KeyWeights map[string]float64

	// MinScore filters out results below this threshold (inclusive).
	MinScore float64
}

// WeightedResult pairs a DiffResult with its computed score.
type WeightedResult struct {
	Result DiffResult
	Score  float64
}

// DefaultWeightOptions returns sensible defaults.
func DefaultWeightOptions() WeightOptions {
	return WeightOptions{
		AddedScore:   1.0,
		RemovedScore: 2.0,
		ChangedScore: 1.5,
		KeyWeights:   nil,
		MinScore:     0,
	}
}

// Weight scores each DiffResult and returns only those meeting MinScore.
func Weight(results []DiffResult, opts WeightOptions) []WeightedResult {
	out := make([]WeightedResult, 0, len(results))
	for _, r := range results {
		var base float64
		switch r.Status {
		case StatusAdded:
			base = opts.AddedScore
		case StatusRemoved:
			base = opts.RemovedScore
		case StatusChanged:
			base = opts.ChangedScore
		default:
			base = 0
		}
		if mul, ok := opts.KeyWeights[r.Key]; ok {
			base *= mul
		}
		if base >= opts.MinScore {
			out = append(out, WeightedResult{Result: r, Score: base})
		}
	}
	return out
}

// HasWeightedResults returns true if any results were produced.
func HasWeightedResults(results []WeightedResult) bool {
	return len(results) > 0
}
