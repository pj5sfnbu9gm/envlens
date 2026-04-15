// Package resolver handles parsing and validation of deployment target
// specifications supplied by the user on the command line.
//
// A target specification has the form:
//
//	name=path
//
// where name is an arbitrary label (e.g. "staging") and path is the
// filesystem path to a .env file.  Paths may contain environment variable
// references (e.g. $HOME) which are expanded before use.  Relative paths
// are resolved against the provided base directory (typically the current
// working directory).
//
// Example:
//
//	targets, err := resolver.ResolveTargets(
//		[]string{"prod=./prod.env", "staging=./staging.env"},
//		"/workspace",
//	)
package resolver
