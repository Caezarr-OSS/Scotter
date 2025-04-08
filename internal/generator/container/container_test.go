package container

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

func TestGenerateContainerFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-container-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test templates directory
	templatesDir := filepath.Join(tempDir, "templates")
	containerTemplatesDir := filepath.Join(templatesDir, "container")
	err = os.MkdirAll(containerTemplatesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}

	// Create a test template file
	testTemplate := `FROM golang:1.21-alpine
WORKDIR /app
COPY . .
CMD ["./{{.ProjectName}}"]`

	err = os.WriteFile(filepath.Join(containerTemplatesDir, "go.dockerfile.tmpl"), []byte(testTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Change to the temp directory for testing
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		err := os.Chdir(originalWd)
		if err != nil {
			t.Logf("Failed to change back to original directory: %v", err)
		}
	}()
	
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test cases
	testCases := []struct {
		name           string
		config         *model.Config
		expectedFile   string
		shouldContain  string
	}{
		{
			name: "Dockerfile for Go project",
			config: &model.Config{
				ProjectName: "testproject",
				Language:    model.GoLang,
				Pipeline: model.PipelineConfig{
					ContainerFormat: model.DockerfileFormat,
				},
			},
			expectedFile:  "Dockerfile",
			shouldContain: "FROM golang:1.21-alpine",
		},
		{
			name: "Containerfile for Go project",
			config: &model.Config{
				ProjectName: "testproject",
				Language:    model.GoLang,
				Pipeline: model.PipelineConfig{
					ContainerFormat: model.ContainerfileFormat,
				},
			},
			expectedFile:  "Containerfile",
			shouldContain: "FROM golang:1.21-alpine",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up any previous test files
			os.Remove(tc.expectedFile)

			// Create generator
			generator := NewGenerator(tc.config, templatesDir)

			// Generate container file
			err := generator.GenerateContainerFile()
			if err != nil {
				t.Fatalf("Failed to generate container file: %v", err)
			}

			// Check if the file was created
			if _, err := os.Stat(tc.expectedFile); os.IsNotExist(err) {
				t.Fatalf("Expected file %s was not created", tc.expectedFile)
			}

			// Read the file content
			content, err := os.ReadFile(tc.expectedFile)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			// Check if the content contains expected string
			if string(content) == "" {
				t.Fatalf("Generated file is empty")
			}

			if tc.shouldContain != "" && !contains(string(content), tc.shouldContain) {
				t.Fatalf("Generated file does not contain expected content.\nExpected to contain: %s\nActual content: %s", 
					tc.shouldContain, string(content))
			}
		})
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return strings.HasPrefix(s, substr) || strings.HasSuffix(s, substr) || strings.Contains(s, substr)
}
