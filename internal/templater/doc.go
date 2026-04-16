// Package templater renders text templates using environment variable maps.
//
// It exposes a small set of template functions:
//
//	{{ env "KEY" }}         — substitutes the value of KEY (empty string if missing)
//	{{ envOr "KEY" "def" }} — substitutes KEY or falls back to a default
//	{{ upper "text" }}      — converts text to upper case
//	{{ lower "text" }}      — converts text to lower case
//
// The caller can choose to fail on missing keys via Options.FailOnMissing.
package templater
