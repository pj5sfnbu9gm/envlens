// Package differ provides utilities for comparing environment variable maps.
//
// The Scope function restricts a diff operation to a subset of keys defined
// by explicit key names or key prefixes. This is useful when you only care
// about a particular service's variables (e.g. DB_, REDIS_) within a larger
// environment file.
package differ
