package differ

import (
	"testing"
)

func TestFingerprint_Empty(t *testing.T) {
	results := Fingerprint(nil, DefaultFingerprintOptions())
	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}

func TestFingerprint_SingleTarget(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"HOST": "example.com", "PORT": "443"},
	}
	results := Fingerprint(targets, DefaultFingerprintOptions())
	if len(results) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(results))
	}
	if results[0].Target != "prod" {
		t.Errorf("unexpected target name: %s", results[0].Target)
	}
	if results[0].KeyCount != 2 {
		t.Errorf("expected KeyCount=2, got %d", results[0].KeyCount)
	}
	if len(results[0].Fingerprint) != 64 {
		t.Errorf("expected 64-char hex fingerprint, got %d chars", len(results[0].Fingerprint))
	}
}

func TestFingerprint_DeterministicOrder(t *testing.T) {
	env := map[string]string{"Z": "1", "A": "2", "M": "3"}
	targets := map[string]map[string]string{"t": env}
	r1 := Fingerprint(targets, DefaultFingerprintOptions())
	r2 := Fingerprint(targets, DefaultFingerprintOptions())
	if r1[0].Fingerprint != r2[0].Fingerprint {
		t.Error("fingerprint is not deterministic")
	}
}

func TestFingerprint_DifferentEnvs_DifferentHashes(t *testing.T) {
	targets := map[string]map[string]string{
		"staging": {"HOST": "staging.example.com"},
		"prod":    {"HOST": "prod.example.com"},
	}
	results := Fingerprint(targets, DefaultFingerprintOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
	if results[0].Fingerprint == results[1].Fingerprint {
		t.Error("expected different fingerprints for different envs")
	}
}

func TestFingerprint_SameEnvs_SameHash(t *testing.T) {
	env := map[string]string{"HOST": "same.example.com", "PORT": "8080"}
	targets := map[string]map[string]string{
		"a": {"HOST": "same.example.com", "PORT": "8080"},
		"b": {"HOST": "same.example.com", "PORT": "8080"},
	}
	_ = env
	results := Fingerprint(targets, DefaultFingerprintOptions())
	if results[0].Fingerprint != results[1].Fingerprint {
		t.Error("expected identical fingerprints for identical envs")
	}
}

func TestFingerprint_KeysOnly_IgnoresValues(t *testing.T) {
	opts := FingerprintOptions{IncludeValues: false}
	targets := map[string]map[string]string{
		"x": {"HOST": "alpha"},
		"y": {"HOST": "beta"},
	}
	results := Fingerprint(targets, opts)
	if results[0].Fingerprint != results[1].Fingerprint {
		t.Error("expected same fingerprint when values are excluded")
	}
}

func TestFingerprint_OnlyKeys_Filters(t *testing.T) {
	opts := FingerprintOptions{IncludeValues: true, OnlyKeys: []string{"HOST"}}
	targets := map[string]map[string]string{
		"t": {"HOST": "h", "SECRET": "s"},
	}
	results := Fingerprint(targets, opts)
	if results[0].KeyCount != 1 {
		t.Errorf("expected KeyCount=1, got %d", results[0].KeyCount)
	}
}

func TestHasFingerprintConflicts_True(t *testing.T) {
	entries := []FingerprintEntry{
		{Target: "a", Fingerprint: "abc123"},
		{Target: "b", Fingerprint: "abc123"},
	}
	if !HasFingerprintConflicts(entries) {
		t.Error("expected conflict to be detected")
	}
}

func TestHasFingerprintConflicts_False(t *testing.T) {
	entries := []FingerprintEntry{
		{Target: "a", Fingerprint: "abc123"},
		{Target: "b", Fingerprint: "def456"},
	}
	if HasFingerprintConflicts(entries) {
		t.Error("expected no conflict")
	}
}
