package github

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Test that the workflow generator can correctly process all templates
func TestAllWorkflowTemplates(t *testing.T) {
	// Create a temporary directory for tests
	tempDir, err := os.MkdirTemp("", "scotter-github-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

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

	// Create the generator
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
