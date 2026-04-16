// Package deduper identifies and removes duplicate environment variable keys.
//
// In a standard Go map there can be no duplicate keys, but when merging
// env files from multiple sources (or when case-sensitivity is a concern)
// logical duplicates can appear. Dedupe provides a CaseFold mode that treats
// keys differing only in case as duplicates, retaining the lexicographically
// first key and recording which keys were dropped.
package deduper
