package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Test that the workflow generator can correctly process all templates
func TestAllWorkflowTemplates(t *testing.T) {
	// In CI environments, templates might not be available, so we'll skip this test
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping template tests in CI environment")
	}

	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "scotter-github-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Setup template directories
	templatesDir := filepath.Join(tempDir, "templates", "github")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates directory: %v", err)
	}

	// Copy template files from project templates directory if available
	// This allows test to run both locally and in CI
	tryTemplatePaths := []string{
		"../../../internal/templates/github",
		"../../internal/templates/github",
		"internal/templates/github",
	}
	
	templateFiles := []string{"ci.yml.tmpl", "commitlint.yml.tmpl", "changelog.yml.tmpl", "release.yml.tmpl"}
	templatesCopied := false
	
	for _, path := range tryTemplatePaths {
		if _, err := os.Stat(path); err == nil {
			// Templates directory found, copy template files
			for _, file := range templateFiles {
				srcPath := filepath.Join(path, file)
				dstPath := filepath.Join(templatesDir, file)
				
				if _, err := os.Stat(srcPath); err == nil {
					// Read source template
					content, err := os.ReadFile(srcPath)
					if err != nil {
						t.Logf("Could not read template %s: %v", srcPath, err)
						continue
					}
					
					// Write to destination
					err = os.WriteFile(dstPath, content, 0644)
					if err != nil {
						t.Logf("Could not write template %s: %v", dstPath, err)
						continue
					}
				}
			}
			templatesCopied = true
			t.Logf("Templates copied from %s to %s", path, templatesDir)
			break
		}
	}
	
	if !templatesCopied {
		// Create minimal template files for testing if none were copied
		for _, file := range templateFiles {
			dstPath := filepath.Join(templatesDir, file)
			// Create a minimal template file with just enough content to parse
			minimalContent := "name: Test Workflow\non: [push]\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v3"
			err = os.WriteFile(dstPath, []byte(minimalContent), 0644)
			if err != nil {
				t.Fatalf("Could not create minimal template %s: %v", dstPath, err)
			}
		}
		t.Log("Created minimal template files for testing")
	}

	// Base configuration for tests
	config := &model.Config{
		ProjectName: "test-project",
		Language:    model.GoLang,
		Go: model.GoConfig{
			ProjectType: model.CLIGoType,
			ModulePath:  "github.com/test/test-project",
		},
		Pipeline: model.PipelineConfig{
			UseGitHubActions: true,
			SelectedFeatures: []string{"ci", "commit-lint", "changelog", "release"},
		},
	}

	// Create the generator - specify the templates directory we created
	generator := NewGenerator(config, tempDir)

	// Test each workflow individually
	t.Run("CI Workflow", func(t *testing.T) {
		if err := generator.GenerateCIWorkflow(); err != nil {
			t.Fatalf("Failed to generate CI workflow: %v", err)
		}
		validateWorkflowFile(t, tempDir, "ci.yml")
	})

	t.Run("Commit Lint Workflow", func(t *testing.T) {
		if err := generator.GenerateCommitLintWorkflow(); err != nil {
			t.Fatalf("Failed to generate commit lint workflow: %v", err)
		}
		validateWorkflowFile(t, tempDir, "commitlint.yml")
	})

	t.Run("Changelog Workflow", func(t *testing.T) {
		if err := generator.GenerateChangelogWorkflow(); err != nil {
			t.Fatalf("Failed to generate changelog workflow: %v", err)
		}
		validateWorkflowFile(t, tempDir, "changelog.yml")
	})

	t.Run("Release Workflow", func(t *testing.T) {
		if err := generator.GenerateReleaseWorkflow(); err != nil {
			t.Fatalf("Failed to generate release workflow: %v", err)
		}
		validateWorkflowFile(t, tempDir, "release.yml")
	})
}

// Verify that a workflow file exists and is not empty
func validateWorkflowFile(t *testing.T, baseDir, filename string) {
	workflowPath := filepath.Join(baseDir, ".github", "workflows", filename)
	
	// Check that the file exists
	info, err := os.Stat(workflowPath)
	if err != nil {
		t.Fatalf("Workflow file %s does not exist: %v", filename, err)
	}
	
	// Check that the file is not empty
	if info.Size() == 0 {
		t.Fatalf("Workflow file %s is empty", filename)
	}
}
