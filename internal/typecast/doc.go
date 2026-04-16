// Package typecast provides utilities for casting environment variable
// string values into typed Go values such as int, float64, and bool.
//
// Usage:
//
//	opts := typecast.DefaultOptions()
//	opts.Hints = map[string]string{
//		"PORT":  "int",
//		"DEBUG": "bool",
//		"RATIO": "float",
//	}
//	results, err := typecast.Cast(env, opts)
//
// Keys not present in Hints are returned as strings.
// In lenient mode (default), cast failures are recorded in Result.Error.
// In strict mode, the first failure returns an error immediately.
package typecast
