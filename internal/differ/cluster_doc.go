// Package differ provides functions for comparing environment variable
// configurations across multiple deployment targets.
//
// Cluster groups keys that share identical non-empty values across two or more
// targets, helping identify redundant configuration, shared secrets, or
// accidental value reuse between environments.
package differ
