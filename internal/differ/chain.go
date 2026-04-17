package differ

// ChainResult holds the diff between consecutive targets in a chain.
type ChainResult struct {
	From    string
	To      string
	Results []Result
}

// Chain computes diffs between each consecutive pair of targets in order.
// targets is an ordered map represented as a slice of named envs.
func Chain(targets []NamedEnv) []ChainResult {
	if len(targets) < 2 {
		return nil
	}
	out := make([]ChainResult, 0, len(targets)-1)
	for i := 1; i < len(targets); i++ {
		prev := targets[i-1]
		curr := targets[i]
		out = append(out, ChainResult{
			From:    prev.Name,
			To:      curr.Name,
			Results: Diff(prev.Env, curr.Env),
		})
	}
	return out
}

// NamedEnv pairs a name with an environment map.
type NamedEnv struct {
	Name string
	Env  map[string]string
}

// HasChainChanges returns true if any step in the chain has changes.
func HasChainChanges(chain []ChainResult) bool {
	for _, c := range chain {
		for _, r := range c.Results {
			if r.Status != StatusUnchanged {
				return true
			}
		}
	}
	return false
}
