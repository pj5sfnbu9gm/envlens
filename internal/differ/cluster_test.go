package differ

import (
	"testing"
)

func TestCluster_Empty(t *testing.T) {
	result := Cluster(nil, DefaultClusterOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d entries", len(result))
	}
}

func TestCluster_SharedValueAcrossTargets(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"DB_HOST": "db.example.com", "CACHE_HOST": "db.example.com"},
		"staging": {"DB_HOST": "db.example.com", "CACHE_HOST": "db.example.com"},
	}
	result := Cluster(targets, DefaultClusterOptions())
	if len(result) == 0 {
		t.Fatal("expected at least one cluster")
	}
	found := false
	for _, e := range result {
		if e.Value == "db.example.com" {
			found = true
			if len(e.Keys) < 2 {
				t.Errorf("expected at least 2 keys, got %d", len(e.Keys))
			}
		}
	}
	if !found {
		t.Error("expected cluster for db.example.com")
	}
}

func TestCluster_MinTargetsFilters(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"X": "shared", "Y": "shared"},
		"staging": {"X": "other"},
	}
	opts := DefaultClusterOptions()
	opts.MinTargets = 2
	result := Cluster(targets, opts)
	for _, e := range result {
		if e.Value == "shared" && len(e.Targets) < 2 {
			t.Errorf("cluster should require %d targets", opts.MinTargets)
		}
	}
}

func TestCluster_MinKeysFilters(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"ONLY_KEY": "v"},
		"staging": {"ONLY_KEY": "v"},
	}
	opts := DefaultClusterOptions()
	opts.MinKeys = 2
	result := Cluster(targets, opts)
	for _, e := range result {
		if len(e.Keys) < opts.MinKeys {
			t.Errorf("cluster should have at least %d keys", opts.MinKeys)
		}
	}
}

func TestCluster_EmptyValuesSkipped(t *testing.T) {
	targets := map[string]map[string]string{
		"prod":    {"A": "", "B": ""},
		"staging": {"A": "", "B": ""},
	}
	result := Cluster(targets, DefaultClusterOptions())
	if len(result) != 0 {
		t.Fatalf("expected no clusters for empty values, got %d", len(result))
	}
}

func TestHasClusters_True(t *testing.T) {
	entries := []ClusterEntry{{Value: "x", Keys: []string{"A"}, Targets: []string{"prod"}}}
	if !HasClusters(entries) {
		t.Error("expected HasClusters to return true")
	}
}

func TestHasClusters_False(t *testing.T) {
	if HasClusters(nil) {
		t.Error("expected HasClusters to return false")
	}
}
