package differ

import "sort"

// TrendDirection indicates whether a metric is increasing, decreasing, or stable.
type TrendDirection string

const (
	TrendUp     TrendDirection = "up"
	TrendDown   TrendDirection = "down"
	TrendStable TrendDirection = "stable"
)

// TrendEntry summarises how a single key's change-count evolved across snapshots.
type TrendEntry struct {
	Key       string
	Counts    []int          // change count per window (oldest first)
	Total     int
	Direction TrendDirection
}

// TrendOptions controls Trend behaviour.
type TrendOptions struct {
	MinTotal        int  // skip keys whose total change count is below this
	IncludeStable   bool // include keys that never changed
}

// DefaultTrendOptions returns sensible defaults.
func DefaultTrendOptions() TrendOptions {
	return TrendOptions{MinTotal: 1, IncludeStable: false}
}

// Trend analyses a slice of MultiDiff results (oldest first) and returns a
// per-key trend showing how frequently each key changed over time.
func Trend(windows []map[string][]DiffResult, opts TrendOptions) []TrendEntry {
	if len(windows) == 0 {
		return nil
	}

	// Collect all keys across all windows.
	keySet := map[string]struct{}{}
	for _, w := range windows {
		for _, results := range w {
			for _, r := range results {
				keySet[r.Key] = struct{}{}
			}
		}
	}

	var entries []TrendEntry
	for key := range keySet {
		counts := make([]int, len(windows))
		for i, w := range windows {
			for _, results := range w {
				for _, r := range results {
					if r.Key == key && r.Status != StatusUnchanged {
						counts[i]++
					}
				}
			}
		}
		total := 0
		for _, c := range counts {
			total += c
		}
		if total < opts.MinTotal && !opts.IncludeStable {
			continue
		}
		entries = append(entries, TrendEntry{
			Key:       key,
			Counts:    counts,
			Total:     total,
			Direction: trendDirection(counts),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Total != entries[j].Total {
			return entries[i].Total > entries[j].Total
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// HasTrendChanges returns true when at least one entry has a non-zero total.
func HasTrendChanges(entries []TrendEntry) bool {
	for _, e := range entries {
		if e.Total > 0 {
			return true
		}
	}
	return false
}

func trendDirection(counts []int) TrendDirection {
	if len(counts) < 2 {
		return TrendStable
	}
	half := len(counts) / 2
	var first, second int
	for i := 0; i < half; i++ {
		first += counts[i]
	}
	for i := half; i < len(counts); i++ {
		second += counts[i]
	}
	switch {
	case second > first:
		return TrendUp
	case second < first:
		return TrendDown
	default:
		return TrendStable
	}
}
