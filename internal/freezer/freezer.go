package freezer

import (
	"errors"
	"fmt"
)

// FrozenEnv holds an immutable snapshot of an environment map.
type FrozenEnv struct {
	data map[string]string
}

// Options controls Freeze behaviour.
type Options struct {
	// AllowEmpty permits keys with empty values.
	AllowEmpty bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{AllowEmpty: true}
}

// Freeze creates a FrozenEnv from the given map.
// A deep copy is taken so subsequent mutations to src have no effect.
func Freeze(src map[string]string, opts Options) (*FrozenEnv, error) {
	if src == nil {
		return nil, errors.New("freezer: source map must not be nil")
	}
	copy := make(map[string]string, len(src))
	for k, v := range src {
		if k == "" {
			return nil, fmt.Errorf("freezer: empty key is not allowed")
		}
		if !opts.AllowEmpty && v == "" {
			return nil, fmt.Errorf("freezer: key %q has empty value and AllowEmpty is false", k)
		}
		copy[k] = v
	}
	return &FrozenEnv{data: copy}, nil
}

// Get returns the value for key and whether it existed.
func (f *FrozenEnv) Get(key string) (string, bool) {
	v, ok := f.data[key]
	return v, ok
}

// Keys returns all keys in the frozen env.
func (f *FrozenEnv) Keys() []string {
	keys := make([]string, 0, len(f.data))
	for k := range f.data {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of entries.
func (f *FrozenEnv) Len() int {
	return len(f.data)
}

// ToMap returns a mutable deep copy of the underlying data.
func (f *FrozenEnv) ToMap() map[string]string {
	out := make(map[string]string, len(f.data))
	for k, v := range f.data {
		out[k] = v
	}
	return out
}

// Diff returns keys whose values differ between two frozen envs.
func Diff(a, b *FrozenEnv) []string {
	seen := make(map[string]struct{})
	var changed []string
	for k, av := range a.data {
		seen[k] = struct{}{}
		if bv, ok := b.data[k]; !ok || bv != av {
			changed = append(changed, k)
		}
	}
	for k := range b.data {
		if _, ok := seen[k]; !ok {
			changed = append(changed, k)
		}
	}
	return changed
}
