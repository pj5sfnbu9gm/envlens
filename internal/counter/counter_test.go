package counter

import (
	"testing"
)

func TestCount_Empty(t *testing.T) {
	s := Count(map[string]string{})
	if s.Total != 0 || s.Empty != 0 || s.NonEmpty != 0 {
		t.Errorf("expected zero stats, got %+v", s)
	}
}

func TestCount_Total(t *testing.T) {
	env := map[string]string{"A": "1", "B": "", "C": "3"}
	s := Count(env)
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
	if s.NonEmpty != 2 {
		t.Errorf("expected NonEmpty=2, got %d", s.NonEmpty)
	}
}

func TestCount_Prefixes(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "prod",
	}
	s := Count(env)
	if s.Prefixes["DB"] != 2 {
		t.Errorf("expected DB prefix count 2, got %d", s.Prefixes["DB"])
	}
	if s.Prefixes["APP"] != 1 {
		t.Errorf("expected APP prefix count 1, got %d", s.Prefixes["APP"])
	}
}

func TestCount_LongestKey(t *testing.T) {
	env := map[string]string{"SHORT": "a", "MUCH_LONGER_KEY": "b"}
	s := Count(env)
	if s.LongestKey != "MUCH_LONGER_KEY" {
		t.Errorf("expected LongestKey=MUCH_LONGER_KEY, got %s", s.LongestKey)
	}
}

func TestTopPrefixes(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "h",
		"DB_PORT":  "p",
		"DB_NAME":  "n",
		"APP_ENV":  "e",
		"APP_PORT": "8080",
		"LOG_LEVEL": "info",
	}
	s := Count(env)
	top := TopPrefixes(s, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 top prefixes, got %d", len(top))
	}
	if top[0] != "DB" {
		t.Errorf("expected first prefix DB, got %s", top[0])
	}
}

func TestTopPrefixes_FewerThanN(t *testing.T) {
	env := map[string]string{"X_ONE": "1"}
	s := Count(env)
	top := TopPrefixes(s, 5)
	if len(top) != 1 {
		t.Errorf("expected 1 prefix, got %d", len(top))
	}
}
