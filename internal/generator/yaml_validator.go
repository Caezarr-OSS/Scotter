package generator

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// ValidateYAML checks if the content of a YAML file is valid
func ValidateYAML(content []byte) error {
	var out interface{}
	err := yaml.Unmarshal(content, &out)
	if err != nil {
		return fmt.Errorf("Invalid YAML: %w", err)
	}
	return nil
}

// ValidateYAMLFile checks if a YAML file is valid
func ValidateYAMLFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}
	return ValidateYAML(content)
}

// ValidateYAMLReader checks if a YAML stream is valid
func ValidateYAMLReader(r io.Reader) error {
	content, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("unable to read stream: %w", err)
	}
	return ValidateYAML(content)
}
