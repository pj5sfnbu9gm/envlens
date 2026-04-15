package differ_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/differ"
)

func TestDiff_Unchanged(t *testing.T) {
	left := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	right := map[string]string{"APP_ENV": "production", "PORT": "8080"}

	result := differ.Diff(left, right)

	if result.HasChanges() {
		t.Fatal("expected no changes, but HasChanges() returned true")
	}
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
	for _, e := range result.Entries {
		if e.Status != differ.StatusUnchanged {
			t.Errorf("key %q: expected unchanged, got %s", e.Key, e.Status)
		}
	}
}

func TestDiff_Added(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{"NEW_KEY": "value"}

	result := differ.Diff(left, right)

	if !result.HasChanges() {
		t.Fatal("expected changes, but HasChanges() returned false")
	}
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Status != differ.StatusAdded {
		t.Errorf("expected added, got %s", result.Entries[0].Status)
	}
	if result.Entries[0].NewValue != "value" {
		t.Errorf("unexpected NewValue: %q", result.Entries[0].NewValue)
	}
}

func TestDiff_Removed(t *testing.T) {
	left := map[string]string{"OLD_KEY": "old"}
	right := map[string]string{}

	result := differ.Diff(left, right)

	if len(result.Entries) != 1 || result.Entries[0].Status != differ.StatusRemoved {
		t.Errorf("expected removed entry, got %+v", result.Entries)
	}
	if result.Entries[0].OldValue != "old" {
		t.Errorf("unexpected OldValue: %q", result.Entries[0].OldValue)
	}
}

func TestDiff_Changed(t *testing.T) {
	left := map[string]string{"DB_URL": "postgres://localhost/dev"}
	right := map[string]string{"DB_URL": "postgres://prod-host/prod"}

	result := differ.Diff(left, right)

	if len(result.Entries) != 1 || result.Entries[0].Status != differ.StatusChanged {
		t.Errorf("expected changed entry, got %+v", result.Entries)
	}
}

func TestDiff_SortedKeys(t *testing.T) {
	left := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	right := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}

	result := differ.Diff(left, right)

	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, e := range result.Entries {
		if e.Key != expected[i] {
			t.Errorf("position %d: expected key %q, got %q", i, expected[i], e.Key)
		}
	}
}

func TestDiff_Mixed(t *testing.T) {
	left := map[string]string{"KEEP": "same", "REMOVE": "gone", "CHANGE": "old"}
	right := map[string]string{"KEEP": "same", "ADD": "new", "CHANGE": "new"}

	result := differ.Diff(left, right)

	counts := map[differ.Status]int{}
	for _, e := range result.Entries {
		counts[e.Status]++
	}

	if counts[differ.StatusAdded] != 1 {
		t.Errorf("expected 1 added, got %d", counts[differ.StatusAdded])
	}
	if counts[differ.StatusRemoved] != 1 {
		t.Errorf("expected 1 removed, got %d", counts[differ.StatusRemoved])
	}
	if counts[differ.StatusChanged] != 1 {
		t.Errorf("expected 1 changed, got %d", counts[differ.StatusChanged])
	}
	if counts[differ.StatusUnchanged] != 1 {
		t.Errorf("expected 1 unchanged, got %d", counts[differ.StatusUnchanged])
	}
}
