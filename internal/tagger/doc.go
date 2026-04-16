// Package tagger assigns user-defined tags to environment variable keys
// based on prefix patterns or explicit key mappings.
//
// Tags are arbitrary string labels (e.g. "sensitive", "database", "infra")
// that can be used downstream for filtering, reporting, or auditing.
//
// Example:
//
//	opts := tagger.DefaultOptions()
//	opts.PrefixTags["DB_"] = []string{"database"}
//	opts.ExplicitTags["SECRET_KEY"] = []string{"sensitive"}
//	results := tagger.Tag(env, opts)
package tagger
