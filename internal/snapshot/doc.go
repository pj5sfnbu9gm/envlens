// Package snapshot provides functionality to persist and restore
// environment variable states to and from disk.
//
// Snapshots are stored as JSON files containing the target name,
// a UTC timestamp, and the full map of environment key-value pairs.
// They can be used to compare current configurations against a
// previously captured baseline.
//
// Example usage:
//
//	// Save current env state
//	err := snapshot.Save("./snapshots/prod.json", "prod", envMap)
//
//	// Load a previously saved snapshot
//	snap, err := snapshot.Load("./snapshots/prod.json")
package snapshot
