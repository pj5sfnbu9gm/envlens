package grouper_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/grouper"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_PASSWORD": "secret",
		"AWS_KEY":     "AKIA123",
		"AWS_SECRET":  "abc",
		"PORT":        "8080",
	}
}

func TestGroup_DefaultDelimiter(t *testing.T) {
	result := grouper.Group(sampleEnv(), grouper.DefaultOptions())
	if len(result["DB"]) != 3 {
		t.Errorf("expected 3 DB keys, got %d", len(result["DB"]))
	}
	if len(result["AWS"]) != 2 {
		t.Errorf("expected 2 AWS keys, got %d", len(result["AWS"]))
	}
	if _, ok := result[""]["PORT"]; !ok {
		t.Error("expected PORT in ungrouped bucket")
	}
}

func TestGroup_MinGroupSize(t *testing.T) {
	opts := grouper.DefaultOptions()
	opts.MinGroupSize = 3
	result := grouper.Group(sampleEnv(), opts)
	if _, ok := result["AWS"]; ok {
		t.Error("AWS group should be filtered out (only 2 members)")
	}
	if _, ok := result["DB"]; !ok {
		t.Error("DB group should be present (3 members)")
	}
}

func TestGroup_EmptyEnv(t *testing.T) {
	result := grouper.Group(map[string]string{}, grouper.DefaultOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d groups", len(result))
	}
}

func TestGroup_CustomDelimiter(t *testing.T) {
	env := map[string]string{
		"APP.HOST": "localhost",
		"APP.PORT": "9000",
		"OTHER":    "val",
	}
	opts := grouper.Options{Delimiter: ".", MinGroupSize: 0}
	result := grouper.Group(env, opts)
	if len(result["APP"]) != 2 {
		t.Errorf("expected 2 APP keys, got %d", len(result["APP"]))
	}
}

func TestGroup_NoDelimiterInAnyKey(t *testing.T) {
	env := map[string]string{"HOST": "a", "PORT": "b"}
	result := grouper.Group(env, grouper.DefaultOptions())
	if len(result[""]) != 2 {
		t.Errorf("expected 2 ungrouped keys, got %d", len(result[""]))
	}
}
