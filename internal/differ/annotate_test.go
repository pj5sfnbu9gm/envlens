package differ

import (
	"testing"
)

func TestAnnotate_DefaultLabels(t *testing.T) {
	results := []Result{
		{Key: "FOO", Status: StatusAdded, NewValue: "bar"},
		{Key: "BAZ", Status: StatusRemoved, OldValue: "old"},
	}
	opts := DefaultAnnotateOptions()
	anns := Annotate(results, opts)
	if len(anns) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(anns))
	}
	if anns[0].Label != "[+]" {
		t.Errorf("expected [+], got %s", anns[0].Label)
	}
	if anns[1].Label != "[-]" {
		t.Errorf("expected [-], got %s", anns[1].Label)
	}
}

func TestAnnotate_KeyNote(t *testing.T) {
	results := []Result{
		{Key: "SECRET", Status: StatusChanged, OldValue: "a", NewValue: "b"},
	}
	opts := DefaultAnnotateOptions()
	opts.KeyNotes["SECRET"] = "sensitive value changed"
	anns := Annotate(results, opts)
	if anns[0].Note != "sensitive value changed" {
		t.Errorf("unexpected note: %s", anns[0].Note)
	}
}

func TestAnnotate_Unchanged(t *testing.T) {
	results := []Result{
		{Key: "STABLE", Status: StatusUnchanged, OldValue: "v", NewValue: "v"},
	}
	opts := DefaultAnnotateOptions()
	anns := Annotate(results, opts)
	if anns[0].Label != "[=]" {
		t.Errorf("expected [=], got %s", anns[0].Label)
	}
}

func TestAnnotate_EmptyResults(t *testing.T) {
	opts := DefaultAnnotateOptions()
	anns := Annotate([]Result{}, opts)
	if len(anns) != 0 {
		t.Errorf("expected empty annotations")
	}
}

func TestAnnotate_MissingNote(t *testing.T) {
	results := []Result{
		{Key: "PLAIN", Status: StatusAdded, NewValue: "x"},
	}
	opts := DefaultAnnotateOptions()
	anns := Annotate(results, opts)
	if anns[0].Note != "" {
		t.Errorf("expected empty note, got %q", anns[0].Note)
	}
}
