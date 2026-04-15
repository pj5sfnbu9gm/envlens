// Package loader provides utilities for reading and parsing .env files
// into structured EnvMap representations used by envlens for auditing
// and diffing environment variable configurations.
//
// Supported .env file format:
//
//	# This is a comment
//	KEY=value
//	QUOTED="some value"
//	SINGLE='another value'
//
// Blank lines and lines beginning with '#' are ignored.
// Values may optionally be wrapped in single or double quotes, which
// will be stripped during parsing.
//
// Example usage:
//
//	env, err := loader.LoadFile(".env.production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(env["APP_ENV"])
package loader
