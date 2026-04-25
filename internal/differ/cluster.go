package differ

import (
	"sort"
)

// ClusterEntry holds a group of keys that share the same value across all targets.
type ClusterEntry struct {
	Value   string
	Keys    []string
	Targets []string
}

// ClusterOptions controls how Cluster behaves.
type ClusterOptions struct {
	// MinTargets is the minimum number of targets a value must appear in.
	MinTargets int
	// MinKeys is the minimum number of keys that must share the value.
	MinKeys int
}

// DefaultClusterOptions returns sensible defaults.
func DefaultClusterOptions() ClusterOptions {
	return ClusterOptions{
		MinTargets: 2,
		MinKeys:    2,
	}
}

// Cluster groups keys across targets that share identical values, surfacing
// potential redundancy or shared configuration.
func Cluster(targets map[string]map[string]string, opts ClusterOptions) []ClusterEntry {
	if len(targets) == 0 {
		return nil
	}

	// value -> target -> []keys
	valueTargetKeys := map[string]map[string][]string{}

	for target, env := range targets {
		for key, val := range env {
			if val == "" {
				continue
			}
			if valueTargetKeys[val] == nil {
				valueTargetKeys[val] = map[string][]string{}
			}
			valueTargetKeys[val][target] = append(valueTargetKeys[val][target], key)
		}
	}

	var results []ClusterEntry
	for val, targetMap := range valueTargetKeys {
		if len(targetMap) < opts.MinTargets {
			continue
		}
		keySet := map[string]struct{}{}
		var tgts []string
		for t, keys := range targetMap {
			tgts = append(tgts, t)
			for _, k := range keys {
				keySet[k] = struct{}{}
			}
		}
		if len(keySet) < opts.MinKeys {
			continue
		}
		var keys []string
		for k := range keySet {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sort.Strings(tgts)
		results = append(results, ClusterEntry{
			Value:   val,
			Keys:    keys,
			Targets: tgts,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value < results[j].Value
	})
	return results
}

// HasClusters returns true if any cluster entries were found.
func HasClusters(entries []ClusterEntry) bool {
	return len(entries) > 0
}
