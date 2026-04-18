// Package differ provides utilities for comparing environment variable maps
// across deployment targets.
//
// The overlap sub-feature identifies keys that are present in some targets
// but absent from others, helping teams detect configuration drift where
// certain variables have not been propagated to all environments.
//
// Use FindOverlap to get a detailed report of which keys are missing from
// which targets, and HasOverlap as a quick boolean check.
package differ
