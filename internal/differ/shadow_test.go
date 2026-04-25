package differ

import (
	"testing"
)

func TestShadow_Empty(t *testing.T) {
	out := Shadow(nil, nil, DefaultShadowOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(out))
	}
}

func TestShadow_NoDifferences(t *testing.T) {
	primary := map[string]string{"APP_ENV": "prod", "PORT": "8080"}
	shadows := map[string]map[string]string{
		"staging": {"APP_ENV": "prod", "PORT": "8080"},
	}
	out := Shadow(primary, shadows, DefaultShadowOptions())
	if len(out) != 0 {
		t.Fatalf("expected no discrepancies, got %d", len(out))
	}
}

func TestShadow_ValueDiffers(t *testing.T) {
	primary := map[string]string{"APP_ENV": "prod"}
	shadows := map[string]map[string]string{
		"staging": {"APP_ENV": "staging"},
	}
	out := Shadow(primary, shadows, DefaultShadowOptions())
	if len(out["APP_ENV"]) != 1 {
		t.Fatalf("expected 1 entry for APP_ENV, got %d", len(out["APP_ENV"]))
	}
	e := out["APP_ENV"][0]
	if e.PrimaryValue != "prod" || e.ShadowValue != "staging" {
		t.Errorf("unexpected values: primary=%q shadow=%q", e.PrimaryValue, e.ShadowValue)
	}
	if e.OnlyInShadow || e.OnlyInPrimary {
		t.Error("expected both sides to have the key")
	}
}

func TestShadow_OnlyInShadow(t *testing.T) {
	primary := map[string]string{}
	shadows := map[string]map[string]string{
		"staging": {"SECRET": "xyz"},
	}
	out := Shadow(primary, shadows, DefaultShadowOptions())
	if len(out["SECRET"]) != 1 {
		t.Fatalf("expected entry for SECRET")
	}
	if !out["SECRET"][0].OnlyInShadow {
		t.Error("expected OnlyInShadow=true")
	}
}

func TestShadow_OnlyInPrimary(t *testing.T) {
	primary := map[string]string{"DB_URL": "postgres://localhost/db"}
	shadows := map[string]map[string]string{
		"staging": {},
	}
	out := Shadow(primary, shadows, DefaultShadowOptions())
	if len(out["DB_URL"]) != 1 {
		t.Fatalf("expected entry for DB_URL")
	}
	if !out["DB_URL"][0].OnlyInPrimary {
		t.Error("expected OnlyInPrimary=true")
	}
}

func TestShadow_IncludeUnchanged(t *testing.T) {
	primary := map[string]string{"PORT": "8080"}
	shadows := map[string]map[string]string{
		"staging": {"PORT": "8080"},
	}
	opts := DefaultShadowOptions()
	opts.IncludeUnchanged = true
	out := Shadow(primary, shadows, opts)
	if len(out["PORT"]) != 1 {
		t.Fatalf("expected PORT entry when IncludeUnchanged=true")
	}
}

func TestHasShadowDifferences(t *testing.T) {
	empty := map[string][]ShadowEntry{}
	if HasShadowDifferences(empty) {
		t.Error("expected false for empty map")
	}
	withEntry := map[string][]ShadowEntry{
		"KEY": {{Key: "KEY"}},
	}
	if !HasShadowDifferences(withEntry) {
		t.Error("expected true when entries present")
	}
}
