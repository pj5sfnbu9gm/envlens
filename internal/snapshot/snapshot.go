package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of environment variables at a point in time.
type Snapshot struct {
	Target    string            `json:"target"`
	Timestamp time.Time         `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// Save writes a snapshot to a JSON file at the given path.
func Save(path string, target string, env map[string]string) error {
	snap := Snapshot{
		Target:    target,
		Timestamp: time.Now().UTC(),
		Env:       env,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	if snap.Target == "" {
		return nil, fmt.Errorf("snapshot: missing target field")
	}

	return &snap, nil
}
