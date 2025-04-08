package github

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

func TestGenerateWorkflows(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Sauvegarder le répertoire de travail actuel
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir) // Revenir au répertoire d'origine à la fin du test
	
	// Changer le répertoire de travail vers le répertoire temporaire
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory to temp dir: %v", err)
	}
	
	t.Logf("Changed working directory to: %s", tempDir)

	// Create necessary subdirectories
	workflowsDir := filepath.Join(".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatalf("Failed to create workflows dir: %v", err)
	}

	// Create a templates directory and copy our template
	templatesDir := filepath.Join("templates", "github")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}

	// Create templates for testing
	changelogTemplate := `name: Generate Changelog
on:
  workflow_dispatch:
jobs:
  generate-changelog:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: echo "Generate changelog"`

	ciTemplate := `name: CI
on:
  push:
    branches: [ main, master ]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: echo "Run tests"`

	// Write the changelog template
	if err := os.WriteFile(filepath.Join(templatesDir, "changelog.yml.tmpl"), []byte(changelogTemplate), 0644); err != nil {
		t.Fatalf("Failed to write changelog template file: %v", err)
	}
	
	// Write the CI template
	if err := os.WriteFile(filepath.Join(templatesDir, "ci.yml.tmpl"), []byte(ciTemplate), 0644); err != nil {
		t.Fatalf("Failed to write CI template file: %v", err)
	}

	// Test cases
	testCases := []struct {
		name      string
		features  []string
		expectFile string
	}{
		{
			name:      "With changelog feature",
			features:  []string{"changelog"},
			expectFile: "changelog.yml",
		},
		{
			name:      "Without changelog feature",
			features:  []string{"ci"},
			expectFile: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean the workflows directory
			files, _ := os.ReadDir(workflowsDir)
			for _, f := range files {
				os.Remove(filepath.Join(workflowsDir, f.Name()))
			}

			// Create a test config with the specified features
			cfg := &model.Config{
				ProjectName: "test-project",
				Pipeline: model.PipelineConfig{
					UseGitHubActions: true,
					SelectedFeatures: tc.features,
				},
			}

			// Create a generator with the test config
			t.Logf("Using templates directory: %s", templatesDir)
			generator := &Generator{
				Config:      cfg,
				TemplateDir: templatesDir,
			}

			// Generate all workflows
			err = generator.Generate()
			if err != nil {
				t.Fatalf("Failed to generate workflows: %v", err)
			}

			// Check if the expected files were created
			if tc.expectFile != "" {
				workflowPath := filepath.Join(workflowsDir, tc.expectFile)
				if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
					t.Errorf("Expected workflow file %s was not created", tc.expectFile)
				} else {
					// Check file content
					content, err := os.ReadFile(workflowPath)
					if err != nil {
						t.Fatalf("Failed to read workflow file: %v", err)
					}
					if !strings.Contains(string(content), "Generate Changelog") {
						t.Errorf("Workflow file does not contain expected content")
					}
				}
			} else {
				// Check that the file was not created
				files, _ := os.ReadDir(workflowsDir)
				for _, f := range files {
					if f.Name() == "changelog.yml" {
						t.Errorf("Changelog workflow file was created when it should not have been")
					}
				}
			}
		})
	}
}

func TestFeatureSelection(t *testing.T) {
	// Test that the correct features are detected
	testCases := []struct {
		name           string
		features       []string
		expectChangelog bool
	}{
		{
			name:           "With changelog feature",
			features:       []string{"changelog"},
			expectChangelog: true,
		},
		{
			name:           "With release feature (which depends on changelog)",
			features:       []string{"release"},
			expectChangelog: true, // Should be true because release depends on changelog
		},
		{
			name:           "Without changelog feature",
			features:       []string{"ci", "commit-lint"},
			expectChangelog: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test config with the specified features
			cfg := &model.Config{
				ProjectName: "test-project",
				Pipeline: model.PipelineConfig{
					UseGitHubActions: true,
					SelectedFeatures: tc.features,
				},
			}

			// Resolve dependencies
			cfg.Pipeline.SelectedFeatures = model.ResolveFeatureDependencies(cfg.Pipeline.SelectedFeatures)

			// Check if changelog is included
			hasChangelog := false
			for _, feature := range cfg.Pipeline.SelectedFeatures {
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
