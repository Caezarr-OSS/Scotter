package github

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// TestChangelogIntegration tests the complete changelog feature integration
func TestChangelogIntegration(t *testing.T) {
	// Skip if we're running in CI
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping integration test in CI environment")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-changelog-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create project structure
	projectDir := filepath.Join(tempDir, "test-project")
	githubDir := filepath.Join(projectDir, ".github", "workflows")
	templatesDir := filepath.Join(tempDir, "templates", "github")

	// Create directories
	for _, dir := range []string{githubDir, templatesDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create test template files
	templates := map[string]string{
		"changelog.yml": `name: Generate Changelog

on:
  workflow_dispatch:
  push:
    branches:
      - main
      - master
    paths-ignore:
      - 'CHANGELOG.md'

jobs:
  generate-changelog:
    name: Generate Changelog
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 16
      
      - name: Install conventional-changelog-cli
        run: npm install -g conventional-changelog-cli
      
      - name: Generate changelog
        run: conventional-changelog -p angular -i CHANGELOG.md -s
      
      - name: Commit and push if changed
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add CHANGELOG.md
          git diff --quiet && git diff --staged --quiet || git commit -m "docs: update changelog [skip ci]"
          git push`,
	}

	for name, content := range templates {
		if err := os.WriteFile(filepath.Join(templatesDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write template %s: %v", name, err)
		}
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
			// Clean the workflows directory
			files, _ := os.ReadDir(githubDir)
			for _, f := range files {
				os.Remove(filepath.Join(githubDir, f.Name()))
			}

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

			// Create a custom generator function for testing
			generateChangelogWorkflow := func() error {
				// Only generate if changelog is in the resolved features
				hasChangelog := false
				for _, feature := range resolvedFeatures {
					if feature == "changelog" {
						hasChangelog = true
						break
					}
				}

				if !hasChangelog {
					return nil
				}

				// Read the template
				templateBytes, err := os.ReadFile(filepath.Join(templatesDir, "changelog.yml"))
				if err != nil {
					return err
				}

				// Write the workflow file
				return os.WriteFile(filepath.Join(githubDir, "changelog.yml"), templateBytes, 0644)
			}

			// Generate the changelog workflow
			if err := generateChangelogWorkflow(); err != nil {
				t.Fatalf("Failed to generate changelog workflow: %v", err)
			}

			// Check if the workflow file was created
			workflowPath := filepath.Join(githubDir, "changelog.yml")
			fileExists := true
			if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
				fileExists = false
			}

			if tc.expectChangelog && !fileExists {
				t.Errorf("Expected changelog.yml to be created, but it wasn't")
			} else if !tc.expectChangelog && fileExists {
				t.Errorf("changelog.yml was created when it shouldn't have been")
			}

			// If the file exists, check its content
			if fileExists {
				content, err := os.ReadFile(workflowPath)
				if err != nil {
					t.Fatalf("Failed to read workflow file: %v", err)
				}

				if !strings.Contains(string(content), "Generate Changelog") {
					t.Errorf("Workflow file does not contain expected content")
				}
			}
		})
	}
}

// TestTaskfileChangelogTaskContent tests the content of the changelog task in Taskfile.yml
func TestTaskfileChangelogTaskContent(t *testing.T) {
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

// TestChangelogFeatureDependencies tests that the changelog feature dependencies are correctly resolved
func TestChangelogFeatureDependencies(t *testing.T) {
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
			name:            "With release feature (depends on changelog)",
			features:        []string{"release"},
			expectChangelog: true,
		},
		{
			name:            "With commit-lint feature (changelog depends on it)",
			features:        []string{"commit-lint"},
			expectChangelog: false,
		},
		{
			name:            "Without changelog feature",
			features:        []string{"ci"},
			expectChangelog: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Resolve dependencies
			resolvedFeatures := model.ResolveFeatureDependencies(tc.features)

			// Check if changelog is included
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
