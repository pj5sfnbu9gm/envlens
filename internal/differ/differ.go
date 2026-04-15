// Package differ provides functionality to compare environment variable
// configurations between two deployment targets, producing a structured diff.
package differ

// Status represents the state of an environment variable key in a diff.
type Status string

const (
	// StatusAdded indicates the key exists only in the right-hand map.
	StatusAdded Status = "added"
	// StatusRemoved indicates the key exists only in the left-hand map.
	StatusRemoved Status = "removed"
	// StatusChanged indicates the key exists in both maps but with different values.
	StatusChanged Status = "changed"
	// StatusUnchanged indicates the key exists in both maps with the same value.
	StatusUnchanged Status = "unchanged"
)

// Entry represents a single key in the diff result.
type Entry struct {
	Key      string
	Status   Status
	OldValue string
	NewValue string
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if the diff contains any added, removed, or changed entries.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Status != StatusUnchanged {
			return true
		}
	}
	return false
}

// Diff compares two environment variable maps (left vs right) and returns a Result.
// Keys are reported in a stable order: sorted alphabetically.
func Diff(left, right map[string]string) *Result {
	seen := make(map[string]bool)
	var entries []Entry

	// Collect all keys from both maps.
	allKeys := make([]string, 0, len(left)+len(right))
	for k := range left {
		if !seen[k] {
			allKeys = append(allKeys, k)
			seen[k] = true
		}
	}
	for k := range right {
		if !seen[k] {
			allKeys = append(allKeys, k)
			seen[k] = true
		}
	}

	// Sort keys for deterministic output.
	sortStrings(allKeys)

	for _, k := range allKeys {
		lv, inLeft := left[k]
		rv, inRight := right[k]

		switch {
		case inLeft && inRight && lv == rv:
			entries = append(entries, Entry{Key: k, Status: StatusUnchanged, OldValue: lv, NewValue: rv})
		case inLeft && inRight:
			entries = append(entries, Entry{Key: k, Status: StatusChanged, OldValue: lv, NewValue: rv})
		case inLeft:
			entries = append(entries, Entry{Key: k, Status: StatusRemoved, OldValue: lv})
		default:
			entries = append(entries, Entry{Key: k, Status: StatusAdded, NewValue: rv})
		}
	}

	return &Result{Entries: entries}
}

// sortStrings sorts a string slice in-place using a simple insertion sort
// to avoid importing "sort" solely for this small helper.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
