// Package freezer provides an immutable snapshot type for environment maps.
//
// A FrozenEnv is created via Freeze and cannot be modified after creation.
// It supports safe concurrent reads, key lookup, and diffing between two
// frozen environments.
//
// Example:
//
//	f, err := freezer.Freeze(env, freezer.DefaultOptions())
//	if err != nil { ... }
//	v, ok := f.Get("DATABASE_URL")
package freezer
