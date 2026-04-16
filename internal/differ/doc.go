// Package differ compares two environment variable maps and produces
// a list of DiffResult entries describing what changed between them.
//
// Each result carries a Status of Added, Removed, Changed, or Unchanged,
// along with the old and new values for the key.
package differ
