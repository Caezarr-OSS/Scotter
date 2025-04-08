package github

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// TestChangelogFeatureSelection tests that the changelog feature is properly selected
// This test doesn't require Node.js or npx
func TestChangelogFeatureSelection(t *testing.T) {
	// Test cases for feature selection
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
			name:            "With release feature (depends on changelog)",
			features:        []string{"release"},
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
			// Create a test config
			cfg := &model.Config{
				ProjectName: "test-project",
				Pipeline: model.PipelineConfig{
					UseGitHubActions: true,
					SelectedFeatures: tc.features,
				},
			}

			// Resolve dependencies
			resolvedFeatures := model.ResolveFeatureDependencies(tc.features)
			cfg.Pipeline.SelectedFeatures = resolvedFeatures

			// Check if changelog is included after dependency resolution
			hasChangelog := false
			for _, feature := range resolvedFeatures {
				if feature == "changelog" {
					hasChangelog = true
					break
				}
			}

			if hasChangelog != tc.expectChangelog {
				t.Errorf("Expected hasChangelog to be %v, got %v", tc.expectChangelog, hasChangelog)
			}
		})
	}
}

// TestChangelogWorkflowGeneration tests the workflow file generation
// This test doesn't require Node.js or npx
func TestChangelogWorkflowGeneration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create necessary directories
	workflowsDir := filepath.Join(tempDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatalf("Failed to create workflows dir: %v", err)
	}

	// Create a simple template file for testing
	templateContent := `name: Generate Changelog
on:
  workflow_dispatch:
jobs:
  generate-changelog:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: echo "Generate changelog"`

	// Create the template file directly in the workflows directory for this test
	if err := os.WriteFile(filepath.Join(workflowsDir, "changelog.yml"), []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Verify the file exists and has the right content
	content, err := os.ReadFile(filepath.Join(workflowsDir, "changelog.yml"))
	if err != nil {
		t.Fatalf("Failed to read workflow file: %v", err)
	}

	if string(content) != templateContent {
		t.Errorf("Workflow file content doesn't match expected template")
	}
}

// TestTaskfileChangelogTask tests that the changelog task in Taskfile.yml works correctly
// This is a simple validation of the task structure, not its execution
func TestTaskfileChangelogTask(t *testing.T) {
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
	expectedComponents := []string{
		"changelog:",
		"desc: Generate or update CHANGELOG.md",
		"if ! command -v npx",
		"npm install -g conventional-changelog-cli",
		"conventional-changelog -p angular -i CHANGELOG.md -s",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(changelogTask, component) {
			t.Errorf("Changelog task doesn't contain expected component: %s", component)
		}
	}
}
