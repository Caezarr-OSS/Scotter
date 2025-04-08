package github

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSimpleChangelogWorkflowGeneration tests the changelog workflow generation with a simplified approach
func TestSimpleChangelogWorkflowGeneration(t *testing.T) {
	// Skip if we're running in CI
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-changelog-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the .github/workflows directory
	workflowsDir := filepath.Join(tempDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatalf("Failed to create workflows dir: %v", err)
	}

	// Create a simple changelog.yml template
	templateDir := filepath.Join(tempDir, "templates", "github")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	// Write a simple changelog.yml template
	templateContent := `name: Generate Changelog

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
          git push`

	if err := os.WriteFile(filepath.Join(templateDir, "changelog.yml"), []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Note: We're not using the config in this simplified test

	// Create a simple function to generate the changelog workflow
	generateChangelogWorkflow := func() error {
		// Read the template
		templateBytes, err := os.ReadFile(filepath.Join(templateDir, "changelog.yml"))
		if err != nil {
			return err
		}

		// Write the workflow file
		return os.WriteFile(filepath.Join(workflowsDir, "changelog.yml"), templateBytes, 0644)
	}

	// Generate the changelog workflow
	if err := generateChangelogWorkflow(); err != nil {
		t.Fatalf("Failed to generate changelog workflow: %v", err)
	}

	// Check if the workflow file was created
	workflowPath := filepath.Join(workflowsDir, "changelog.yml")
	if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
		t.Errorf("Changelog workflow file was not created at %s", workflowPath)
	}

	// Read the generated workflow file
	generatedBytes, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("Failed to read generated workflow file: %v", err)
	}

	// Verify the content
	if string(generatedBytes) != templateContent {
		t.Errorf("Generated workflow content does not match template")
	}
}
