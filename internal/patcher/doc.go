// Package patcher applies a sequence of declarative patch operations
// (set, unset, rename) to an environment variable map, producing a new
// map without mutating the original.
//
// # Patch Operations
//
//   - OpSet   – create or overwrite a key with a given value.
//   - OpUnset – remove a key from the map.
//   - OpRename – move a key to a new name, preserving its value.
//
// Each operation produces a Result that records whether it was applied
// and a human-readable note, making it easy to surface changes in
// higher-level reporters.
package patcher
