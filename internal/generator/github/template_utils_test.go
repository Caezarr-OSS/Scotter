package github

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplateDelimiters(t *testing.T) {
	// Test default delimiters
	defaultDelims := DefaultDelimiters()
	if defaultDelims.Left != "{{" || defaultDelims.Right != "}}" {
		t.Errorf("Default delimiters should be {{ and }}, got %s and %s", defaultDelims.Left, defaultDelims.Right)
	}

	// Test GitHub workflow delimiters
	workflowDelims := GitHubWorkflowDelimiters()
	if workflowDelims.Left != "[[" || workflowDelims.Right != "]]" {
		t.Errorf("GitHub workflow delimiters should be [[ and ]], got %s and %s", workflowDelims.Left, workflowDelims.Right)
	}
}

func TestExecuteTemplateWithDelimiters(t *testing.T) {
	// Create a temporary test template with GitHub workflow delimiters
	tempDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test template with GitHub workflow delimiters
	workflowTemplate := `name: CI
on:
  push:
    branches: [develop, main, [[.BranchPattern]]]
jobs:
  build:
    name: [[.JobName]]
    runs-on: ubuntu-latest`

	templatePath := filepath.Join(tempDir, "workflow.yml.tmpl")
	if err := os.WriteFile(templatePath, []byte(workflowTemplate), 0644); err != nil {
		t.Fatalf("Failed to write test template: %v", err)
	}

	// Test data
	data := struct {
		BranchPattern string
		JobName       string
	}{
		BranchPattern: "feature/*",
		JobName:       "Build and Test",
	}

	// Execute template with GitHub workflow delimiters
	result, err := ExecuteTemplateWithDelimiters(templatePath, data, GitHubWorkflowDelimiters())
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}

	// Check if the template was correctly rendered
	expectedOutput := `name: CI
on:
  push:
    branches: [develop, main, feature/*]
jobs:
  build:
    name: Build and Test
    runs-on: ubuntu-latest`

	if result != expectedOutput {
		t.Errorf("Template output does not match expected value.\nGot:\n%s\n\nExpected:\n%s", result, expectedOutput)
	}

	// Test with invalid template
	invalidTemplate := `name: CI
on:
  push:
    branches: [[.MissingClose`

	invalidPath := filepath.Join(tempDir, "invalid.yml.tmpl")
	if err := os.WriteFile(invalidPath, []byte(invalidTemplate), 0644); err != nil {
		t.Fatalf("Failed to write invalid template: %v", err)
	}

	// Validate should fail
	if err := ValidateTemplate(invalidPath, GitHubWorkflowDelimiters()); err == nil {
		t.Error("ValidateTemplate should fail with invalid template")
	}

	// ExecuteTemplateWithDelimiters should also fail
	_, err = ExecuteTemplateWithDelimiters(invalidPath, data, GitHubWorkflowDelimiters())
	if err == nil {
		t.Error("ExecuteTemplateWithDelimiters should fail with invalid template")
	}
}
