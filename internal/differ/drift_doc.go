// Package differ provides utilities for comparing environment variable maps
// across targets and snapshots.
//
// The drift sub-feature (DetectDrift) analyses an ordered sequence of env
// snapshots for a single named target and surfaces every key that was added,
// removed, or changed between consecutive snapshots.  This is useful for
// tracking configuration churn over time — for example, comparing weekly
// snapshots of a production environment to understand which variables are
// stable and which are frequently mutated.
//
// Usage:
//
//	snaps := []map[string]string{week1, week2, week3}
//	report := differ.DetectDrift("prod", snaps)
//	if differ.HasDrift(report) {
//		for _, e := range report.Entries {
//			fmt.Println(e.Key, e.Status)
//		}
//	}
package differ
