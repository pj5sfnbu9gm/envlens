package differ

// IntersectResult holds a key that exists in all provided envs with its values per target.
type IntersectResult struct {
	Key    string
	Values map[string]string // target name -> value
}

// Intersect returns keys that appear in every target env map.
// targets is a map of target-name -> env map.
func Intersect(targets map[string]map[string]string) []IntersectResult {
	if len(targets) == 0 {
		return nil
	}

	// collect all keys per target
	keyCounts := map[string]int{}
	for _, env := range targets {
		for k := range env {
			keyCounts[k]++
		}
	}

	n := len(targets)
	var results []IntersectResult
	for _, k := range sortStrings(keysOf(keyCounts)) {
		if keyCounts[k] < n {
			continue
		}
		values := make(map[string]string, n)
		for name, env := range targets {
			values[name] = env[k]
		}
		results = append(results, IntersectResult{Key: k, Values: values})
	}
	return results
}

// AllAgree returns true when every target has the same value for the given key.
func AllAgree(results []IntersectResult) []string {
	var agreed []string
	for _, r := range results {
		if allSame(r.Values) {
			agreed = append(agreed, r.Key)
		}
	}
	return agreed
}

func allSame(m map[string]string) bool {
	var ref string
	first := true
	for _, v := range m {
		if first {
			ref = v
			first = false
			continue
		}
		if v != ref {
			return false
		}
	}
	return true
}

func keysOf(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
