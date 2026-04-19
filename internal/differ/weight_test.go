package differ

import (
	"testing"
)

func buildWeightResults() []DiffResult {
	return []DiffResult{
		{Key: "HOST", Status: StatusAdded, NewValue: "localhost"},
		{Key: "PORT", Status: StatusRemoved, OldValue: "8080"},
		{Key: "DEBUG", Status: StatusChanged, OldValue: "false", NewValue: "true"},
		{Key: "STABLE", Status: StatusUnchanged, OldValue: "ok", NewValue: "ok"},
	}
}

func TestWeight_DefaultScores(t *testing.T) {
	opts := DefaultWeightOptions()
	res := Weight(buildWeightResults(), opts)
	// Unchanged should be excluded (score 0, MinScore 0 means >= 0 passes — adjust: unchanged base is 0)
	// With MinScore 0, unchanged (score 0) is included; let's verify counts
	if len(res) != 4 {
		t.Fatalf("expected 4 results, got %d", len(res))
	}
}

func TestWeight_AddedScore(t *testing.T) {
	opts := DefaultWeightOptions()
	res := Weight([]DiffResult{{Key: "A", Status: StatusAdded}}, opts)
	if res[0].Score != 1.0 {
		t.Errorf("expected 1.0, got %f", res[0].Score)
	}
}

func TestWeight_RemovedScore(t *testing.T) {
	opts := DefaultWeightOptions()
	res := Weight([]DiffResult{{Key: "B", Status: StatusRemoved}}, opts)
	if res[0].Score != 2.0 {
		t.Errorf("expected 2.0, got %f", res[0].Score)
	}
}

func TestWeight_ChangedScore(t *testing.T) {
	opts := DefaultWeightOptions()
	res := Weight([]DiffResult{{Key: "C", Status: StatusChanged}}, opts)
	if res[0].Score != 1.5 {
		t.Errorf("expected 1.5, got %f", res[0].Score)
	}
}

func TestWeight_KeyMultiplier(t *testing.T) {
	opts := DefaultWeightOptions()
	opts.KeyWeights = map[string]float64{"SECRET": 3.0}
	res := Weight([]DiffResult{{Key: "SECRET", Status: StatusChanged}}, opts)
	if res[0].Score != 1.5*3.0 {
		t.Errorf("expected 4.5, got %f", res[0].Score)
	}
}

func TestWeight_MinScoreFilters(t *testing.T) {
	opts := DefaultWeightOptions()
	opts.MinScore = 2.0
	res := Weight(buildWeightResults(), opts)
	for _, r := range res {
		if r.Score < 2.0 {
			t.Errorf("result %s has score %f below MinScore", r.Result.Key, r.Score)
		}
	}
}

func TestHasWeightedResults_True(t *testing.T) {
	if !HasWeightedResults([]WeightedResult{{Score: 1}}) {
		t.Error("expected true")
	}
}

func TestHasWeightedResults_False(t *testing.T) {
	if HasWeightedResults(nil) {
		t.Error("expected false")
	}
}
