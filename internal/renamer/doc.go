// Package renamer provides utilities for renaming environment variable keys
// within an env map.
//
// It supports two complementary strategies:
//
//  1. Explicit rules – a list of (From, To) pairs that rename individual keys.
//  2. Prefix replacement – bulk-rename all keys whose names start with a given
//     prefix by substituting that prefix with a new one.
//
// The original map is never mutated; Rename always returns a new map.
//
// Example:
//
//	opts := renamer.Options{
//	    Rules:     []renamer.Rule{{From: "DB_HOST", To: "DATABASE_HOST"}},
//	    OldPrefix: "APP_",
//	    NewPrefix: "SERVICE_",
//	}
//	renamed, err := renamer.Rename(env, opts)
package renamer
