package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a set of environment variables as key-value pairs.
type EnvMap map[string]string

// LoadFile reads a .env file and returns an EnvMap.
// It skips blank lines and lines starting with '#'.
func LoadFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("loader: invalid syntax at %s:%d: %q", path, lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		if key == "" {
			return nil, fmt.Errorf("loader: empty key at %s:%d", path, lineNum)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: error reading %q: %w", path, err)
	}

	return env, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
