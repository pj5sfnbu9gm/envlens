// Package differ — gradient.go
//
// Gradient analyses how individual environment variable values shift across an
// ordered sequence of named targets (e.g. dev → staging → prod).
//
// Unlike Trend, which operates on pre-computed diff windows, Gradient works
// directly on raw env maps and is therefore suitable for snapshot-free
// pipeline comparisons where targets are known up-front.
//
// Each GradientEntry records the per-step values and a change count, making
// it easy to surface which keys are "volatile" (changing at every step) vs
// "shifting" (changing at some steps) vs "stable" (never changing).
package differ
