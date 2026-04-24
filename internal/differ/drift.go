package differ

import "sort"

// DriftEntry represents a key that has drifted (changed) across two or more
// sequential snapshots for a single target.
type DriftEntry struct {
	Key    string
	From   string
	To     string
	Status string // "added", "removed", "changed"
}

// DriftReport holds the drift analysis for a single target across snapshots.
type DriftReport struct {
	Target  string
	Entries []DriftEntry
}

// DetectDrift compares an ordered slice of env snapshots (oldest first) for a
// named target and returns a DriftReport describing every key that changed
// between consecutive snapshots.
func DetectDrift(target string, snapshots []map[string]string) DriftReport {
	report := DriftReport{Target: target}
	if len(snapshots) < 2 {
		return report
	}

	seen := map[string]struct{}{}
	for _, snap := range snapshots {
		for k := range snap {
			seen[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for i := 1; i < len(snapshots); i++ {
			prev := snapshots[i-1]
			curr := snapshots[i]
			prevVal, prevOK := prev[k]
			currVal, currOK := curr[k]

			switch {
			case !prevOK && currOK:
				report.Entries = append(report.Entries, DriftEntry{Key: k, From: "", To: currVal, Status: "added"})
			case prevOK && !currOK:
				report.Entries = append(report.Entries, DriftEntry{Key: k, From: prevVal, To: "", Status: "removed"})
			case prevOK && currOK && prevVal != currVal:
				report.Entries = append(report.Entries, DriftEntry{Key: k, From: prevVal, To: currVal, Status: "changed"})
			}
		}
	}

	return report
}

// HasDrift returns true if the report contains any drift entries.
func HasDrift(r DriftReport) bool {
	return len(r.Entries) > 0
}
