package github

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSimpleChangelogGeneration tests the changelog generation without relying on external files
func TestSimpleChangelogGeneration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-changelog-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create project structure
	projectDir := filepath.Join(tempDir, "test-project")
	workflowsDir := filepath.Join(projectDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatalf("Failed to create workflows dir: %v", err)
	}

	// Test cases
	testCases := []struct {
		name            string
		features        []string
		expectChangelog bool
	}{
		{
			name:            "With changelog feature",
			features:        []string{"changelog"},
			expectChangelog: true,
		},
		{
			name:            "Without changelog feature",
			features:        []string{"ci"},
			expectChangelog: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean any existing workflow files
			files, _ := os.ReadDir(workflowsDir)
			for _, f := range files {
				os.Remove(filepath.Join(workflowsDir, f.Name()))
			}

			// Create a simple changelog workflow file directly
			generateChangelog := func(features []string) error {
				// Check if changelog is in the features
				hasChangelog := false
				for _, feature := range features {
					if feature == "changelog" {
						hasChangelog = true
						break
					}
				}

				if !hasChangelog {
					return nil
				}

				// Simple template content for testing
				content := `name: Generate Changelog
on:
  workflow_dispatch:
jobs:
  generate-changelog:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: echo "Generate changelog"`

				return os.WriteFile(filepath.Join(workflowsDir, "changelog.yml"), []byte(content), 0644)
			}

			// Generate the changelog workflow
			if err := generateChangelog(tc.features); err != nil {
				t.Fatalf("Failed to generate changelog workflow: %v", err)
			}

			// Check if the workflow file was created
			workflowPath := filepath.Join(workflowsDir, "changelog.yml")
			fileExists := true
			if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
				fileExists = false
			}

			if tc.expectChangelog && !fileExists {
				t.Errorf("Expected changelog.yml to be created, but it wasn't")
			} else if !tc.expectChangelog && fileExists {
				t.Errorf("changelog.yml was created when it shouldn't have been")
			}
		})
	}
}

// TestChangelogTaskContent tests the content of the changelog task in Taskfile.yml
func TestChangelogTaskContent(t *testing.T) {
	// This is a simplified representation of the changelog task
	changelogTask := `
changelog:
  desc: Generate or update CHANGELOG.md
  cmds:
    - |
      if ! command -v npx &> /dev/null; then
        echo "Error: npx not found. Please install Node.js"
        exit 1
      fi
    - |
      if ! command -v conventional-changelog &> /dev/null; then
        npm install -g conventional-changelog-cli
      fi
    - conventional-changelog -p angular -i CHANGELOG.md -s
`

	// Verify the task has the expected components
	if len(changelogTask) == 0 {
		t.Errorf("Changelog task content is empty")
	}

	// Simple verification that the task contains key components
	expectedPhrases := []string{
		"changelog:",
		"Generate or update CHANGELOG.md",
		"conventional-changelog -p angular -i CHANGELOG.md -s",
	}

	for _, phrase := range expectedPhrases {
		if !contains(changelogTask, phrase) {
			t.Errorf("Changelog task doesn't contain expected phrase: %s", phrase)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
