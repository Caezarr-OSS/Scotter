package github

import (
	"os"
	"path/filepath"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// MockGenerator is a mock implementation of GitHubGenerator for testing
type MockGenerator struct {
	Config           *model.Config
	OutputDir        string
	GenerateCalled   bool
	CIWorkflowCalled bool
	CommitLintCalled bool
	ChangelogCalled  bool
	ReleaseCalled    bool
	DependabotCalled bool
}

// NewMockGenerator creates a new mock generator
func NewMockGenerator(cfg *model.Config, outputDir string) *MockGenerator {
	return &MockGenerator{
		Config:    cfg,
		OutputDir: outputDir,
	}
}

// Generate implements GitHubGenerator.Generate
func (m *MockGenerator) Generate() error {
	m.GenerateCalled = true

	// Check which features are enabled
	hasCI := false
	hasCommitLint := false
	hasChangelog := false
	hasRelease := false
	hasDependabot := false

	for _, feature := range m.Config.Pipeline.SelectedFeatures {
		switch feature {
		case "ci":
			hasCI = true
		case "commit-lint":
			hasCommitLint = true
		case "changelog":
			hasChangelog = true
		case "release":
			hasRelease = true
		case "dependabot":
			hasDependabot = true
		}
	}

	// Generate workflows based on selected features
	if hasCI {
		if err := m.GenerateCIWorkflow(); err != nil {
			return err
		}
	}

	if hasCommitLint {
		if err := m.GenerateCommitLintWorkflow(); err != nil {
			return err
		}
	}

	if hasChangelog {
		if err := m.GenerateChangelogWorkflow(); err != nil {
			return err
		}
	}

	if hasRelease {
		if err := m.GenerateReleaseWorkflow(); err != nil {
			return err
		}
	}

	if hasDependabot {
		if err := m.GenerateDependabotConfig(); err != nil {
			return err
		}
	}

	return nil
}

// GenerateCIWorkflow implements GitHubGenerator.GenerateCIWorkflow
func (m *MockGenerator) GenerateCIWorkflow() error {
	m.CIWorkflowCalled = true
	return os.WriteFile(filepath.Join(m.OutputDir, "ci.yml"), []byte("CI Workflow"), 0644)
}

// GenerateCommitLintWorkflow implements GitHubGenerator.GenerateCommitLintWorkflow
func (m *MockGenerator) GenerateCommitLintWorkflow() error {
	m.CommitLintCalled = true
	return os.WriteFile(filepath.Join(m.OutputDir, "commit-lint.yml"), []byte("Commit Lint Workflow"), 0644)
}

// GenerateChangelogWorkflow implements GitHubGenerator.GenerateChangelogWorkflow
func (m *MockGenerator) GenerateChangelogWorkflow() error {
	m.ChangelogCalled = true
	return os.WriteFile(filepath.Join(m.OutputDir, "changelog.yml"), []byte("Changelog Workflow"), 0644)
}

// GenerateReleaseWorkflow implements GitHubGenerator.GenerateReleaseWorkflow
func (m *MockGenerator) GenerateReleaseWorkflow() error {
	m.ReleaseCalled = true
	return os.WriteFile(filepath.Join(m.OutputDir, "release.yml"), []byte("Release Workflow"), 0644)
}

// GenerateDependabotConfig implements GitHubGenerator.GenerateDependabotConfig
func (m *MockGenerator) GenerateDependabotConfig() error {
	m.DependabotCalled = true
	return os.WriteFile(filepath.Join(m.OutputDir, "dependabot.yml"), []byte("Dependabot Config"), 0644)
}
