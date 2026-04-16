package tagger

import "strings"

// Tag represents a label attached to an env key.
type Tag struct {
	Key   string
	Value string
	Tags  []string
}

// Options configures tagging behaviour.
type Options struct {
	// PrefixTags maps a key prefix to a list of tags.
	PrefixTags map[string][]string
	// ExplicitTags maps exact keys to a list of tags.
	ExplicitTags map[string][]string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		PrefixTags:   map[string][]string{},
		ExplicitTags: map[string][]string{},
	}
}

// Tag applies tags to every key in env and returns a slice of Tag results.
func Tag(env map[string]string, opts Options) []Tag {
	results := make([]Tag, 0, len(env))
	for k, v := range env {
		tags := collectTags(k, opts)
		results = append(results, Tag{Key: k, Value: v, Tags: tags})
	}
	sortTags(results)
	return results
}

func collectTags(key string, opts Options) []string {
	seen := map[string]struct{}{}
	var out []string
	add := func(t string) {
		if _, ok := seen[t]; !ok {
			seen[t] = struct{}{}
			out = append(out, t)
		}
	}
	for prefix, tags := range opts.PrefixTags {
		if strings.HasPrefix(key, prefix) {
			for _, t := range tags {
				add(t)
			}
		}
	}
	for _, t := range opts.ExplicitTags[key] {
		add(t)
	}
	return out
}

func sortTags(tags []Tag) {
	for i := 1; i < len(tags); i++ {
		for j := i; j > 0 && tags[j].Key < tags[j-1].Key; j-- {
			tags[j], tags[j-1] = tags[j-1], tags[j]
		}
	}
}
