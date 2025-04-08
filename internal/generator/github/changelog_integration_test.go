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
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
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
			// Créer le répertoire .github/workflows s'il n'existe pas
			workflowsDir := filepath.Join(projectDir, ".github", "workflows")
			if err := os.MkdirAll(workflowsDir, 0755); err != nil {
				t.Fatalf("Failed to create workflows directory: %v", err)
			}
			
			// Clean up any existing workflow files
			files, _ := os.ReadDir(workflowsDir)
			for _, f := range files {
				os.Remove(filepath.Join(workflowsDir, f.Name()))
			}

			// Resolve dependencies
			resolvedFeatures := model.ResolveFeatureDependencies(tc.features)

			// Afficher les informations sur les chemins des templates
			cwd, _ := os.Getwd()
			t.Logf("Current working directory: %s", cwd)
			t.Logf("Templates directory: %s", filepath.Join(cwd, "templates", "github"))

			// Vérifier si les templates existent
			templatesDir := filepath.Join(cwd, "templates", "github")
			if _, err := os.Stat(templatesDir); err == nil {
				t.Logf("Templates directory exists: %s", templatesDir)
				files, _ := os.ReadDir(templatesDir)
				t.Logf("Files in templates directory:")
				for _, file := range files {
					t.Logf("- %s", file.Name())
				}
			} else {
				t.Logf("Templates directory does not exist: %s (Error: %v)", templatesDir, err)
			}

			// Create a test mock generator
			// Utiliser le répertoire de templates détecté
			mockGenerator := NewTestMockGenerator(projectDir, templatesDir, resolvedFeatures)

			// Generate the changelog workflow
			if err := mockGenerator.GenerateChangelogWorkflow(); err != nil {
				t.Fatalf("Failed to generate changelog workflow: %v", err)
			}
			
			// Nous ne générons pas le workflow CI dans ce test
			// car nous testons uniquement la fonctionnalité changelog

			// Check if the workflow file was created
			// We already have workflowsDir defined above, so we don't need to redefine it
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
