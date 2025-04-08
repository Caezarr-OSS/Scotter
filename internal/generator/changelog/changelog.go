package changelog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Generator handles the creation of CHANGELOG and conventional commits config
type Generator struct {
	Config *model.Config
}

// NewGenerator creates a new changelog generator
func NewGenerator(cfg *model.Config) *Generator {
	return &Generator{
		Config: cfg,
	}
}

// Generate creates CHANGELOG.md and conventional commits configuration
// This is a legacy method for backward compatibility
func (g *Generator) Generate() error {
	// Check if any changelog features are selected in the pipeline
	isChangelogEnabled := false
	isCommitLintEnabled := false

	for _, feature := range g.Config.Pipeline.SelectedFeatures {
		if feature == "changelog" {
			isChangelogEnabled = true
		}
		if feature == "commit-lint" {
			isCommitLintEnabled = true
		}
	}

	if !isChangelogEnabled && !isCommitLintEnabled {
		fmt.Println("Skipping changelog and conventional commits setup...")
		return nil
	}

	fmt.Println("Setting up changelog and conventional commits...")

	// Generate CHANGELOG.md if enabled
	if isChangelogEnabled {
		if err := g.GenerateChangelog(); err != nil {
			return fmt.Errorf("failed to generate CHANGELOG.md: %w", err)
		}
	}

	// Generate .commitlintrc.js if conventional commits are enabled
	if isCommitLintEnabled {
		if err := g.GenerateCommitLintConfig(); err != nil {
			return fmt.Errorf("failed to generate commitlint config: %w", err)
		}
	}

	fmt.Println("Changelog and conventional commits setup completed!")
	return nil
}

// GenerateChangelog creates a CHANGELOG.md file
func (g *Generator) GenerateChangelog() error {
	year := time.Now().Year()
	
	content := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A

## [0.1.0] - %d-MM-DD

### Added
- Initial release

[Unreleased]: https://github.com/username/%s/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/username/%s/releases/tag/v0.1.0
`, year, g.Config.ProjectName, g.Config.ProjectName)

	return os.WriteFile("CHANGELOG.md", []byte(content), 0644)
}

// GenerateCommitLintConfig creates a commitlint configuration file
func (g *Generator) GenerateCommitLintConfig() error {
	content := `module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'body-leading-blank': [1, 'always'],
    'body-max-line-length': [2, 'always', 100],
    'footer-leading-blank': [1, 'always'],
    'footer-max-line-length': [2, 'always', 100],
    'header-max-length': [2, 'always', 100],
    'subject-case': [
      2,
      'never',
      ['sentence-case', 'start-case', 'pascal-case', 'upper-case'],
    ],
    'subject-empty': [2, 'never'],
    'subject-full-stop': [2, 'never', '.'],
    'type-case': [2, 'always', 'lower-case'],
    'type-empty': [2, 'never'],
    'type-enum': [
      2,
      'always',
      [
        'build',
        'chore',
        'ci',
        'docs',
        'feat',
        'fix',
        'perf',
        'refactor',
        'revert',
        'style',
        'test',
      ],
    ],
  },
};`

	return os.WriteFile(".commitlintrc.js", []byte(content), 0644)
}

// GenerateCommitMsgHook creates a Git pre-commit hook for commitlint
func (g *Generator) GenerateCommitMsgHook() error {
	// Check if commit-lint is enabled in the pipeline
	isCommitLintEnabled := false
	for _, feature := range g.Config.Pipeline.SelectedFeatures {
		if feature == "commit-lint" {
			isCommitLintEnabled = true
			break
		}
	}

	if !isCommitLintEnabled {
		return nil
	}

	// Ensure .git/hooks directory exists
	hooksDir := filepath.Join(".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Create commit-msg hook
	hookContent := `#!/bin/sh
# Verifies commit messages follow conventional commits format
# Requires commitlint to be installed

if command -v npx &> /dev/null; then
    # If npx is available
    npx --no-install commitlint --edit "$1"
elif command -v commitlint &> /dev/null; then
    # If commitlint is available directly
    commitlint --edit "$1"
else
    echo "warning: commitlint not found, skipping commit message verification"
    exit 0
fi
`

	hookPath := filepath.Join(hooksDir, "commit-msg")
	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to create commit-msg hook: %w", err)
	}

	return nil
}
