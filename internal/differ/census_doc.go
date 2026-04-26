// Package differ provides utilities for comparing environment variable maps
// across deployment targets.
//
// Census counts how many targets define each key, exposing coverage gaps
// where a key is missing from one or more targets. This is useful for
// identifying configuration drift or incomplete rollouts across environments.
package differ
