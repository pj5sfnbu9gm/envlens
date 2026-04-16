package scorer

import (
	"testing"
)

func TestScoreEnv_Perfect(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	opts := DefaultOptions()
	s := ScoreEnv(env, opts)
	if s.FinalScore != 100 {
		t.Errorf("expected 100, got %d", s.FinalScore)
	}
	if len(s.Penalties) != 0 {
		t.Errorf("expected no penalties, got %d", len(s.Penalties))
	}
}

func TestScoreEnv_EmptyValue(t *testing.T) {
	env := map[string]string{"HOST": ""}
	opts := DefaultOptions()
	s := ScoreEnv(env, opts)
	if s.FinalScore != 95 {
		t.Errorf("expected 95, got %d", s.FinalScore)
	}
}

func TestScoreEnv_LowercaseKey(t *testing.T) {
	env := map[string]string{"host": "localhost"}
	opts := DefaultOptions()
	s := ScoreEnv(env, opts)
	if s.FinalScore != 97 {
		t.Errorf("expected 97, got %d", s.FinalScore)
	}
}

func TestScoreEnv_WhitespaceKey(t *testing.T) {
	env := map[string]string{"MY KEY": "val"}
	opts := DefaultOptions()
	s := ScoreEnv(env, opts)
	// whitespace(4) + lowercase check: "MY KEY" != "MY KEY" toUpper is same, so only whitespace
	if s.FinalScore != 96 {
		t.Errorf("expected 96, got %d", s.FinalScore)
	}
}

func TestScoreEnv_FloorAtZero(t *testing.T) {
	env := map[string]string{}
	for i := 0; i < 30; i++ {
		env[string(rune('a'+i%26))+"key"] = ""
	}
	opts := DefaultOptions()
	s := ScoreEnv(env, opts)
	if s.FinalScore < 0 {
		t.Errorf("score should not be negative, got %d", s.FinalScore)
	}
}

func TestScoreEnv_MultiplePenalties(t *testing.T) {
	env := map[string]string{
		"good_KEY": "",
	}
	opts := DefaultOptions()
	s := ScoreEnv(env, opts)
	// empty(5) + lowercase(3) = 8 deducted
	if s.FinalScore != 92 {
		t.Errorf("expected 92, got %d", s.FinalScore)
	}
	if len(s.Penalties) != 2 {
		t.Errorf("expected 2 penalties, got %d", len(s.Penalties))
	}
}
