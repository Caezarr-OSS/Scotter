package github

import (
	"os"
	"path/filepath"
)

// TestMockGenerator is a simplified version of the GitHub generator for testing
type TestMockGenerator struct {
	ProjectDir   string
	TemplatesDir string
	Features     []string
}

// NewTestMockGenerator creates a new mock generator for testing
func NewTestMockGenerator(projectDir, templatesDir string, features []string) *TestMockGenerator {
	return &TestMockGenerator{
		ProjectDir:   projectDir,
		TemplatesDir: templatesDir,
		Features:     features,
	}
}

// GenerateChangelogWorkflow generates a changelog workflow file for testing
func (g *TestMockGenerator) GenerateChangelogWorkflow() error {
	// Check if changelog feature is enabled
	hasChangelog := false
	for _, feature := range g.Features {
		if feature == "changelog" {
			hasChangelog = true
			break
		}
	}

	if !hasChangelog {
		return nil
	}

	// Create the workflows directory if it doesn't exist
	workflowsDir := filepath.Join(g.ProjectDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		return err
	}

	// Template content for testing
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

	// Write the workflow file
	return os.WriteFile(filepath.Join(workflowsDir, "changelog.yml"), []byte(templateContent), 0644)
}

// GenerateCIWorkflow generates a CI workflow file for testing
func (g *TestMockGenerator) GenerateCIWorkflow() error {
	// Check if CI feature is enabled
	hasCI := false
	for _, feature := range g.Features {
		if feature == "ci" {
			hasCI = true
			break
		}
	}

	// Afficher les fonctionnalités activées
	// Vérification des fonctionnalités activées

	if !hasCI {
		return nil
	}

	// Create the workflows directory if it doesn't exist
	workflowsDir := filepath.Join(g.ProjectDir, ".github", "workflows")
	// Vérifier si le répertoire existe déjà
	if _, err := os.Stat(workflowsDir); os.IsNotExist(err) {
		// Créer le répertoire
		if err := os.MkdirAll(workflowsDir, 0755); err != nil {
			return err
		}
	}

	// Template content for testing
	templateContent := `name: CI

on:
  push:
    branches:
      - main
      - master
  pull_request:
    branches:
      - main
      - master

jobs:
  test:
    name: Run Tests
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test ./...`

	// Write the workflow file
	ciPath := filepath.Join(workflowsDir, "ci.yml")
	if err := os.WriteFile(ciPath, []byte(templateContent), 0644); err != nil {
		return err
	}
	return nil
}
