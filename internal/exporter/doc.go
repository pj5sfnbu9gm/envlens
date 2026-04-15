// Package exporter provides utilities for serialising environment variable
// maps back to text files in several common formats.
//
// Supported formats:
//
//   - dotenv  – plain KEY=VALUE pairs, quoting only when necessary
//   - shell   – KEY="VALUE" pairs suitable for sourcing in bash/sh
//   - export  – like shell but prefixed with the `export` keyword
//
// Basic usage:
//
//	opts := exporter.DefaultOptions()
//	err  := exporter.Export(os.Stdout, env, opts)
//
// To write directly to a file:
//
//	err := exporter.ExportToFile(".env.production", env, opts)
package exporter
