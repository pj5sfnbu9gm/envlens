package differ

import (
	"testing"
)

func TestDigest_Empty(t *testing.T) {
	results := Digest(map[string]map[string]string{}, DefaultDigestOptions())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestDigest_SingleTarget(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"APP_ENV": "production", "PORT": "8080"},
	}
	results := Digest(targets, DefaultDigestOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Target != "prod" {
		t.Errorf("unexpected target name: %s", results[0].Target)
	}
	if len(results[0].Digest) != 64 {
		t.Errorf("expected 64-char hex digest, got len=%d", len(results[0].Digest))
	}
}

func TestDigest_DeterministicOrder(t *testing.T) {
	env1 := map[string]string{"B": "2", "A": "1"}
	env2 := map[string]string{"A": "1", "B": "2"}

	r1 := Digest(map[string]map[string]string{"x": env1}, DefaultDigestOptions())
	r2 := Digest(map[string]map[string]string{"x": env2}, DefaultDigestOptions())

	if r1[0].Digest != r2[0].Digest {
		t.Errorf("digests should be equal regardless of map order")
	}
}

func TestDigest_DifferentEnvs_DifferentDigests(t *testing.T) {
	targets := map[string]map[string]string{
		"staging": {"APP_ENV": "staging"},
		"prod":    {"APP_ENV": "production"},
	}
	results := Digest(targets, DefaultDigestOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 results")
	}
	if results[0].Digest == results[1].Digest {
		t.Errorf("different envs should produce different digests")
	}
}

func TestDigest_IncludeKeys(t *testing.T) {
	targets := map[string]map[string]string{
		"a": {"APP_ENV": "prod", "SECRET": "abc"},
		"b": {"APP_ENV": "prod", "SECRET": "xyz"},
	}
	opts := DigestOptions{IncludeKeys: []string{"APP_ENV"}}
	results := Digest(targets, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 results")
	}
	if results[0].Digest != results[1].Digest {
		t.Errorf("digests should match when only shared key is included")
	}
}

func TestHasDigestConflicts_True(t *testing.T) {
	results := []DigestResult{
		{Target: "a", Digest: "aaa"},
		{Target: "b", Digest: "bbb"},
	}
	if !HasDigestConflicts(results) {
		t.Error("expected conflict")
	}
}

func TestHasDigestConflicts_False(t *testing.T) {
	results := []DigestResult{
		{Target: "a", Digest: "aaa"},
		{Target: "b", Digest: "aaa"},
	}
	if HasDigestConflicts(results) {
		t.Error("expected no conflict")
	}
}

func TestHasDigestConflicts_Single(t *testing.T) {
	results := []DigestResult{{Target: "a", Digest: "aaa"}}
	if HasDigestConflicts(results) {
		t.Error("single result should never conflict")
	}
}
