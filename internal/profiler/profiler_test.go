package profiler

import (
	"testing"
)

func TestAnalyze_EmptyMap(t *testing.T) {
	p := Analyze(map[string]string{})
	if p.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", p.TotalKeys)
	}
	if p.EmptyValues != 0 {
		t.Errorf("expected 0 empty values, got %d", p.EmptyValues)
	}
	if p.SensitiveKeys != 0 {
		t.Errorf("expected 0 sensitive keys, got %d", p.SensitiveKeys)
	}
}

func TestAnalyze_TotalKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux", "QUUX": ""}
	p := Analyze(env)
	if p.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", p.TotalKeys)
	}
}

func TestAnalyze_EmptyValues(t *testing.T) {
	env := map[string]string{
		"FOO":   "bar",
		"EMPTY": "",
		"BLANK": "   ",
	}
	p := Analyze(env)
	if p.EmptyValues != 2 {
		t.Errorf("expected 2 empty values, got %d", p.EmptyValues)
	}
}

func TestAnalyze_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD":  "secret",
		"API_KEY":      "key123",
		"APP_NAME":     "myapp",
		"AUTH_TOKEN":   "tok",
		"NORMAL_VALUE": "hello",
	}
	p := Analyze(env)
	if p.SensitiveKeys != 3 {
		t.Errorf("expected 3 sensitive keys, got %d", p.SensitiveKeys)
	}
}

func TestAnalyze_PrefixCounts(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "prod",
	}
	p := Analyze(env)
	if p.PrefixCounts["DB"] != 2 {
		t.Errorf("expected DB prefix count 2, got %d", p.PrefixCounts["DB"])
	}
	if p.PrefixCounts["APP"] != 1 {
		t.Errorf("expected APP prefix count 1, got %d", p.PrefixCounts["APP"])
	}
}

func TestAnalyze_TopPrefixes(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "h", "DB_PORT": "p", "DB_NAME": "n",
		"APP_ENV": "e", "APP_PORT": "p2",
		"SVC_URL": "u",
	}
	p := Analyze(env)
	if len(p.TopPrefixes) == 0 {
		t.Fatal("expected non-empty top prefixes")
	}
	if p.TopPrefixes[0] != "DB" {
		t.Errorf("expected top prefix DB, got %s", p.TopPrefixes[0])
	}
}

func TestAnalyze_NoPrefixKeys(t *testing.T) {
	env := map[string]string{"NOPREFIX": "val"}
	p := Analyze(env)
	if len(p.PrefixCounts) != 0 {
		t.Errorf("expected no prefix counts for single-segment key")
	}
}
