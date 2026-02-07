// Package envfile provides utilities for loading environment variables from
// dotenv file chains (.env, .env.local, .env.{name}, .env.{name}.local).
package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Load reads a dotenv file chain from the given directory and returns a merged
// map of key-value pairs. Files are loaded in order (later files override earlier):
//
//  1. .env
//  2. .env.local
//  3. .env.{envName}         (if envName is non-empty)
//  4. .env.{envName}.local   (if envName is non-empty)
//
// Missing files are silently skipped. Returns an error only if a file exists but
// cannot be read or contains invalid syntax.
func Load(dir string, envName string) (map[string]string, error) {
	result := make(map[string]string)

	files := []string{
		".env",
		".env.local",
	}

	if envName != "" {
		files = append(files, fmt.Sprintf(".env.%s", envName))
		files = append(files, fmt.Sprintf(".env.%s.local", envName))
	}

	for _, filename := range files {
		path := filepath.Join(dir, filename)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Skip missing files
			}
			return nil, fmt.Errorf("failed to read %s: %w", path, err)
		}

		if err := parseEnvFile(data, result); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", path, err)
		}
	}

	return result, nil
}

// parseEnvFile parses KEY=value lines from raw file bytes into the provided map.
// Supports comments (#), empty lines, optional quoting, and export prefixes.
func parseEnvFile(data []byte, vars map[string]string) error {
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional "export " prefix
		line = strings.TrimPrefix(line, "export ")

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("line %d: invalid format (expected KEY=value): %s", i+1, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		vars[key] = value
	}
	return nil
}
