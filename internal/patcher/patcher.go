package patcher

import "fmt"

// Op represents the type of patch operation.
type Op string

const (
	OpSet    Op = "set"
	OpUnset  Op = "unset"
	OpRename Op = "rename"
)

// Patch describes a single mutation to apply to an env map.
type Patch struct {
	Op    Op
	Key   string
	Value string // used by OpSet
	To    string // used by OpRename
}

// Options controls Patch behaviour.
type Options struct {
	// FailOnMissing causes Apply to return an error when a key targeted by
	// OpUnset or OpRename does not exist in the env map.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{FailOnMissing: false}
}

// Result holds the outcome of a single patch operation.
type Result struct {
	Patch   Patch
	Applied bool
	Note    string
}

// Apply executes a slice of Patch operations against env, returning a new map
// and a per-operation result log.
func Apply(env map[string]string, patches []Patch, opts Options) (map[string]string, []Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	results := make([]Result, 0, len(patches))

	for _, p := range patches {
		var res Result
		res.Patch = p

		switch p.Op {
		case OpSet:
			out[p.Key] = p.Value
			res.Applied = true
			res.Note = "set"

		case OpUnset:
			if _, ok := out[p.Key]; !ok {
				if opts.FailOnMissing {
					return nil, nil, fmt.Errorf("patcher: unset: key %q not found", p.Key)
				}
				res.Note = "key not found, skipped"
			} else {
				delete(out, p.Key)
				res.Applied = true
				res.Note = "unset"
			}

		case OpRename:
			val, ok := out[p.Key]
			if !ok {
				if opts.FailOnMissing {
					return nil, nil, fmt.Errorf("patcher: rename: key %q not found", p.Key)
				}
				res.Note = "key not found, skipped"
			} else {
				delete(out, p.Key)
				out[p.To] = val
				res.Applied = true
				res.Note = fmt.Sprintf("renamed to %s", p.To)
			}

		default:
			return nil, nil, fmt.Errorf("patcher: unknown op %q", p.Op)
		}

		results = append(results, res)
	}

	return out, results, nil
}
