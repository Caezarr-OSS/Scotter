package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// WriteFile writes content to a file with specified permissions
// If the file is a YAML file, it validates the content first
func WriteFile(path string, content []byte, perm os.FileMode) error {
	// If the file ends with .yml or .yaml, validate the YAML
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".yml" || ext == ".yaml" {
		// Validate the YAML
		var out interface{}
		err := yaml.Unmarshal(content, &out)
		if err != nil {
			return fmt.Errorf("Invalid YAML in %s: %w", path, err)
		}
	}

	// Create parent directories if needed
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	// Write the file
	return os.WriteFile(path, content, perm)
}
