package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Target represents a named deployment target with an associated env file path.
type Target struct {
	Name string
	Path string
}

// ResolveTargets parses a slice of "name=path" strings into Target values,
// expanding environment variables and resolving relative paths.
func ResolveTargets(specs []string, baseDir string) ([]Target, error) {
	if len(specs) == 0 {
		return nil, fmt.Errorf("resolver: at least one target must be specified")
	}

	seen := make(map[string]struct{})
	targets := make([]Target, 0, len(specs))

	for _, spec := range specs {
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("resolver: invalid target spec %q: expected name=path", spec)
		}

		name := strings.TrimSpace(parts[0])
		rawPath := os.ExpandEnv(strings.TrimSpace(parts[1]))

		if _, dup := seen[name]; dup {
			return nil, fmt.Errorf("resolver: duplicate target name %q", name)
		}
		seen[name] = struct{}{}

		path := rawPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(baseDir, path)
		}

		targets = append(targets, Target{Name: name, Path: path})
	}

	return targets, nil
}

// ValidatePaths checks that every target's file path exists and is a regular file.
func ValidatePaths(targets []Target) error {
	for _, t := range targets {
		info, err := os.Stat(t.Path)
		if err != nil {
			return fmt.Errorf("resolver: target %q: %w", t.Name, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("resolver: target %q: path %q is not a regular file", t.Name, t.Path)
		}
	}
	return nil
}
