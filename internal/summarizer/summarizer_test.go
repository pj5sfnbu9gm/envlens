package summarizer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/summarizer"
)

func TestSummarize_Empty(t *testing.T) {
	s := summarizer.Summarize(map[string]string{}, 5)
	if s.TotalKeys != 0 {
		t.Fatalf("expected 0 keys, got %d", s.TotalKeys)
	}
}

func TestSummarize_TotalKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	s := summarizer.Summarize(env, 5)
	if s.TotalKeys != 3 {
		t.Fatalf("expected 3, got %d", s.TotalKeys)
	}
}

func TestSummarize_EmptyValues(t *testing.T) {
	env := map[string]string{"A": "", "B": "val", "C": ""}
	s := summarizer.Summarize(env, 5)
	if s.EmptyValues != 2 {
		t.Fatalf("expected 2 empty, got %d", s.EmptyValues)
	}
}

func TestSummarize_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"API_TOKEN":   "tok",
		"APP_NAME":    "myapp",
	}
	s := summarizer.Summarize(env, 5)
	if s.SensitiveKeys != 2 {
		t.Fatalf("expected 2 sensitive, got %d", s.SensitiveKeys)
	}
}

func TestSummarize_UniqueValues(t *testing.T) {
	env := map[string]string{"A": "x", "B": "x", "C": "y"}
	s := summarizer.Summarize(env, 5)
	if s.UniqueValues != 2 {
		t.Fatalf("expected 2 unique values, got %d", s.UniqueValues)
	}
}

func TestSummarize_TopPrefixes(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "h", "DB_PORT": "p", "DB_NAME": "n",
		"APP_ENV": "prod",
	}
	s := summarizer.Summarize(env, 1)
	if len(s.TopPrefixes) != 1 {
		t.Fatalf("expected 1 top prefix, got %d", len(s.TopPrefixes))
	}
	if s.TopPrefixes[0].Prefix != "DB" {
		t.Fatalf("expected DB prefix, got %s", s.TopPrefixes[0].Prefix)
	}
	if s.TopPrefixes[0].Count != 3 {
		t.Fatalf("expected count 3, got %d", s.TopPrefixes[0].Count)
	}
}

func TestSummarize_NoPrefixKeys(t *testing.T) {
	env := map[string]string{"HOSTNAME": "h", "PORT": "8080"}
	s := summarizer.Summarize(env, 5)
	if len(s.TopPrefixes) != 0 {
		t.Fatalf("expected no prefixes, got %d", len(s.TopPrefixes))
	}
}
