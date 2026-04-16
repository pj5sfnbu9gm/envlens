package pivotter_test

import (
	"testing"

	"github.com/user/envlens/internal/pivotter"
)

var targets = map[string]map[string]string{
	"prod": {"HOST": "prod.example.com", "PORT": "443", "DEBUG": "false"},
	"staging": {"HOST": "staging.example.com", "PORT": "443", "DEBUG": "true"},
	"dev": {"HOST": "localhost", "PORT": "8080", "DEBUG": "true"},
}

func TestPivot_ExcludesUnchanged(t *testing.T) {
	rows := pivotter.Pivot(targets, pivotter.DefaultOptions())
	for _, r := range rows {
		if r.Key == "PORT" {
			t.Errorf("PORT is identical across all targets and should be excluded")
		}
	}
	if len(rows) == 0 {
		t.Fatal("expected at least one differing row")
	}
}

func TestPivot_IncludesUnchanged(t *testing.T) {
	opts := pivotter.Options{IncludeUnchanged: true}
	rows := pivotter.Pivot(targets, opts)
	found := false
	for _, r := range rows {
		if r.Key == "PORT" {
			found = true
		}
	}
	if !found {
		t.Error("PORT should appear when IncludeUnchanged is true")
	}
}

func TestPivot_RowValues(t *testing.T) {
	rows := pivotter.Pivot(targets, pivotter.Options{IncludeUnchanged: true})
	for _, r := range rows {
		if r.Key == "HOST" {
			if r.Targets["prod"] != "prod.example.com" {
				t.Errorf("unexpected prod HOST: %s", r.Targets["prod"])
			}
			if r.Targets["dev"] != "localhost" {
				t.Errorf("unexpected dev HOST: %s", r.Targets["dev"])
			}
			return
		}
	}
	t.Error("HOST row not found")
}

func TestPivot_MissingKeyInTarget(t *testing.T) {
	t2 := map[string]map[string]string{
		"a": {"FOO": "1", "BAR": "2"},
		"b": {"FOO": "1"},
	}
	rows := pivotter.Pivot(t2, pivotter.Options{IncludeUnchanged: true})
	for _, r := range rows {
		if r.Key == "BAR" {
			if r.Targets["b"] != "" {
				t.Errorf("expected empty string for missing key in target b")
			}
			return
		}
	}
	t.Error("BAR row not found")
}

func TestPivot_EmptyTargets(t *testing.T) {
	rows := pivotter.Pivot(map[string]map[string]string{}, pivotter.DefaultOptions())
	if len(rows) != 0 {
		t.Errorf("expected no rows for empty targets, got %d", len(rows))
	}
}

func TestPivot_SortedRows(t *testing.T) {
	rows := pivotter.Pivot(targets, pivotter.Options{IncludeUnchanged: true})
	for i := 1; i < len(rows); i++ {
		if rows[i].Key < rows[i-1].Key {
			t.Errorf("rows not sorted: %s before %s", rows[i-1].Key, rows[i].Key)
		}
	}
}
