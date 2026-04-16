package flattener_test

import (
	"testing"

	"github.com/user/envlens/internal/flattener"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP__DB__HOST": "localhost",
		"APP__DB__PORT": "5432",
		"PLAIN_KEY":     "value",
		"EMPTY__KEY":    "",
	}
}

func TestFlatten_DefaultSeparator(t *testing.T) {
	results := flattener.Flatten(sampleEnv(), flattener.DefaultOptions())
	m := flattener.ToMap(results)

	if m["APP.DB.HOST"] != "localhost" {
		t.Errorf("expected APP.DB.HOST=localhost, got %q", m["APP.DB.HOST"])
	}
	if m["APP.DB.PORT"] != "5432" {
		t.Errorf("expected APP.DB.PORT=5432, got %q", m["APP.DB.PORT"])
	}
	if m["PLAIN_KEY"] != "value" {
		t.Errorf("expected PLAIN_KEY=value, got %q", m["PLAIN_KEY"])
	}
}

func TestFlatten_SkipEmpty(t *testing.T) {
	opts := flattener.DefaultOptions()
	opts.SkipEmpty = true
	results := flattener.Flatten(sampleEnv(), opts)
	m := flattener.ToMap(results)

	if _, ok := m["EMPTY.KEY"]; ok {
		t.Error("expected EMPTY.KEY to be skipped")
	}
	if len(m) != 3 {
		t.Errorf("expected 3 results, got %d", len(m))
	}
}

func TestFlatten_CustomSeparators(t *testing.T) {
	env := map[string]string{"DB-HOST": "pg", "DB-PORT": "5432"}
	opts := flattener.Options{Separator: "-", OutputSeparator: "/"}
	m := flattener.ToMap(flattener.Flatten(env, opts))

	if m["DB/HOST"] != "pg" {
		t.Errorf("expected DB/HOST=pg, got %q", m["DB/HOST"])
	}
}

func TestFlatten_NoSeparatorInKey(t *testing.T) {
	env := map[string]string{"SIMPLE": "val"}
	m := flattener.ToMap(flattener.Flatten(env, flattener.DefaultOptions()))
	if m["SIMPLE"] != "val" {
		t.Errorf("expected SIMPLE=val, got %q", m["SIMPLE"])
	}
}

func TestFlatten_EmptyEnv(t *testing.T) {
	results := flattener.Flatten(map[string]string{}, flattener.DefaultOptions())
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestToMap_PreservesValues(t *testing.T) {
	results := []flattener.FlattenResult{
		{OriginalKey: "A__B", FlatKey: "A.B", Value: "x"},
		{OriginalKey: "C", FlatKey: "C", Value: "y"},
	}
	m := flattener.ToMap(results)
	if m["A.B"] != "x" || m["C"] != "y" {
		t.Errorf("unexpected map contents: %v", m)
	}
}
