package inspector

import (
	"testing"
)

func makeEnv() map[string]string {
	return map[string]string{
		"APP_NAME":     "envlens",
		"PORT":         "8080",
		"RATIO":        "3.14",
		"DEBUG":        "true",
		"API_KEY":      "s3cr3t",
		"EMPTY_VAR":    "",
		"HOMEPAGE_URL": "https://example.com",
	}
}

func findEntry(entries []Entry, key string) (Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return Entry{}, false
}

func TestInspect_TypeguessString(t *testing.T) {
	e, ok := findEntry(Inspect(makeEnv()), "APP_NAME")
	if !ok {
		t.Fatal("entry not found")
	}
	if e.Typeguess != "string" {
		t.Errorf("expected string, got %s", e.Typeguess)
	}
}

func TestInspect_TypeguessInt(t *testing.T) {
	e, _ := findEntry(Inspect(makeEnv()), "PORT")
	if e.Typeguess != "int" {
		t.Errorf("expected int, got %s", e.Typeguess)
	}
}

func TestInspect_TypeguessFloat(t *testing.T) {
	e, _ := findEntry(Inspect(makeEnv()), "RATIO")
	if e.Typeguess != "float" {
		t.Errorf("expected float, got %s", e.Typeguess)
	}
}

func TestInspect_TypeguessBool(t *testing.T) {
	e, _ := findEntry(Inspect(makeEnv()), "DEBUG")
	if e.Typeguess != "bool" {
		t.Errorf("expected bool, got %s", e.Typeguess)
	}
}

func TestInspect_TypeguessURL(t *testing.T) {
	e, _ := findEntry(Inspect(makeEnv()), "HOMEPAGE_URL")
	if e.Typeguess != "url" {
		t.Errorf("expected url, got %s", e.Typeguess)
	}
}

func TestInspect_EmptyVar(t *testing.T) {
	e, _ := findEntry(Inspect(makeEnv()), "EMPTY_VAR")
	if !e.Empty {
		t.Error("expected Empty=true")
	}
	if e.Typeguess != "empty" {
		t.Errorf("expected empty typeguess, got %s", e.Typeguess)
	}
}

func TestInspect_SensitiveKey(t *testing.T) {
	e, _ := findEntry(Inspect(makeEnv()), "API_KEY")
	if !e.Sensitive {
		t.Error("expected Sensitive=true for API_KEY")
	}
}

func TestInspect_NonSensitiveKey(t *testing.T) {
	e, _ := findEntry(Inspect(makeEnv()), "APP_NAME")
	if e.Sensitive {
		t.Error("expected Sensitive=false for APP_NAME")
	}
}
