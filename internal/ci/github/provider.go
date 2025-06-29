// Package github implements the GitHub Actions CI provider
package github

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caezarr-oss/scotter/internal/embedded"
	"github.com/caezarr-oss/scotter/pkg/plugin"
)

// GitHubProvider implements the CIProvider interface for GitHub Actions
type GitHubProvider struct {
	templateManager plugin.TemplateManager
}

// NewGitHubProvider creates a new GitHub Actions provider
func NewGitHubProvider() *GitHubProvider {
	return &GitHubProvider{
		templateManager: embedded.NewTemplateManager(),
	}
}

// Name returns the CI provider name
func (p *GitHubProvider) Name() string {
	return "github"
}

// SupportedLanguages returns languages supported by this provider
func (p *GitHubProvider) SupportedLanguages() []string {
	return []string{"go"}
}

// GenerateWorkflows generates CI workflows for a language and project type
func (p *GitHubProvider) GenerateWorkflows(projectPath, language, projectType string, config map[string]interface{}) error {
	// Validate language
	if !contains(p.SupportedLanguages(), language) {
		return fmt.Errorf("language '%s' is not supported by GitHub Actions provider", language)
	}

	// Create .github/workflows directory
	workflowsDir := filepath.Join(projectPath, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflows directory: %w", err)
	}

	// Generate CI workflow file
	ciWorkflowPath := filepath.Join(workflowsDir, "ci.yml")
	ciWorkflowContent := generateCIWorkflow(language, projectType)
	if err := os.WriteFile(ciWorkflowPath, []byte(ciWorkflowContent), 0644); err != nil {
		return fmt.Errorf("failed to create CI workflow: %w", err)
	}

	// Generate Release workflow file
	releaseWorkflowPath := filepath.Join(workflowsDir, "release.yml")
	releaseWorkflowContent := generateReleaseWorkflow(language, projectType)
	if err := os.WriteFile(releaseWorkflowPath, []byte(releaseWorkflowContent), 0644); err != nil {
		return fmt.Errorf("failed to create Release workflow: %w", err)
	}

	// Generate Commitlint workflow
	commitlintWorkflowPath := filepath.Join(workflowsDir, "commitlint.yml")
	commitlintWorkflowContent := generateCommitlintWorkflow()
	if err := os.WriteFile(commitlintWorkflowPath, []byte(commitlintWorkflowContent), 0644); err != nil {
		return fmt.Errorf("failed to create Commitlint workflow: %w", err)
	}

	// Create commitlint.config.js in the project root
	commitlintConfigPath := filepath.Join(projectPath, "commitlint.config.js")
	commitlintConfigContent := `module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'body-max-line-length': [1, 'always', 100],
  },
};
`
	if err := os.WriteFile(commitlintConfigPath, []byte(commitlintConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to create commitlint config: %w", err)
	}

	return nil
}

// generateCIWorkflow creates the content for the CI workflow
func generateCIWorkflow(language, projectType string) string {
	// Using fixes identified in memories - avoiding conditional syntax issues
	return `name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: [1.21.x]

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
        
    - name: Build
      run: go build -v ./...
      
    - name: Test
      run: go test -v ./...
`
}

// generateReleaseWorkflow creates the content for the Release workflow
// Incorporates fixes from memories - using RELEASE_TOKEN and git for changelogs
// Adapts tag format based on project type: 'v*' for libraries, all tags for CLI/API/default
func generateReleaseWorkflow(language, projectType string) string {
	// Determine the tag pattern based on project type
	tagPattern := "*"
	if projectType == "library" {
		tagPattern = "'v*'" // Libraries must use 'v' prefix
	} else {
		tagPattern = "'*'" // CLI/API/default can use any SemVer format
	}

	return `name: Release

on:
  push:
    tags:
      - ` + tagPattern + `

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
          
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21.x
          
      - name: Verify tests pass
        run: go test -v ./...

      - name: Install Syft
        run: |
          curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin
          syft --version
          
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v4
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          # Use RELEASE_TOKEN instead of GITHUB_TOKEN as per identified fix
          GITHUB_TOKEN: ${{ secrets.RELEASE_TOKEN }}
`
}

// generateCommitlintWorkflow creates the content for the Commitlint workflow
// Now checks both push events and pull request events to ensure commit standards
func generateCommitlintWorkflow() string {
	return `name: Commitlint

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  lint-commits:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
          
      - uses: wagoid/commitlint-github-action@v5
`
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
