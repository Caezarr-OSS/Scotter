package model

import (
	"fmt"
)

// ProjectType defines the type of Go project to create
type ProjectType string

const (
	// DefaultType is a minimal Go project structure
	DefaultType ProjectType = "default"
	// LibraryType is for Go packages/libraries
	LibraryType ProjectType = "library"
	// CLIType is for command-line applications
	CLIType ProjectType = "cli"
	// APIType is for HTTP API/service applications
	APIType ProjectType = "api"
	// CompleteType includes all features
	CompleteType ProjectType = "complete"
)

// String returns the string representation of ProjectType
func (pt ProjectType) String() string {
	switch pt {
	case DefaultType:
		return "Default"
	case LibraryType:
		return "Library"
	case CLIType:
		return "CLI"
	case APIType:
		return "API"
	case CompleteType:
		return "Complete"
	default:
		return "Unknown"
	}
}

// GitHubFeatures holds GitHub related project features
type GitHubFeatures struct {
	// UseWorkflows specifies whether to include GitHub Actions workflows
	UseWorkflows bool
	// UseCommitLint specifies whether to include commitlint configuration
	UseCommitLint bool
	// UseReleaseWorkflow specifies whether to set up automatic releases
	UseReleaseWorkflow bool
	// UseDependabot specifies whether to set up Dependabot
	UseDependabot bool
	// GenerateChangelog specifies whether to generate CHANGELOG.md template
	GenerateChangelog bool
}

// Features holds all project features
type Features struct {
	// UseTaskFile specifies whether to include a Taskfile
	UseTaskFile bool
	// UseMakeFile specifies whether to include a Makefile
	UseMakeFile bool
	// GitHub features
	GitHub GitHubFeatures
}

// Config holds the project configuration
type Config struct {
	// ProjectName is the name of the project
	ProjectName string
	// ProjectType determines the structure and features
	ProjectType ProjectType
	// ModulePath is the Go module path (e.g., "github.com/username/project")
	ModulePath string
	// Features to include in the project
	Features Features
	// Directories to create
	Directories []string
}

// Validate validates that the configuration is complete and correct
func (cfg *Config) Validate() error {
	// Check required fields
	if cfg.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check project type
	validType := false
	for _, pt := range []ProjectType{DefaultType, LibraryType, CLIType, APIType, CompleteType} {
		if cfg.ProjectType == pt {
			validType = true
			break
		}
	}

	if !validType {
		return fmt.Errorf("invalid project type: %s", cfg.ProjectType)
	}

	return nil
}

// DefaultFeaturesForType returns the default features for a project type
func DefaultFeaturesForType(projectType ProjectType) Features {
	features := Features{}
	
	switch projectType {
	case DefaultType:
		// Minimal features for default type
		features.UseTaskFile = false
		features.GitHub.UseWorkflows = false
	case LibraryType:
		// Library typically has documentation and CI
		features.UseTaskFile = false
		features.GitHub.UseWorkflows = true
		features.GitHub.GenerateChangelog = true
	case CLIType:
		// CLI apps typically need build tools and releases
		features.UseTaskFile = true
		features.GitHub.UseWorkflows = true
		features.GitHub.UseReleaseWorkflow = true
	case APIType:
		// API services need build and deployment tools
		features.UseTaskFile = true
		features.GitHub.UseWorkflows = true
	case CompleteType:
		// Complete includes all features
		features.UseTaskFile = true
		features.GitHub.UseWorkflows = true
		features.GitHub.UseCommitLint = true
		features.GitHub.UseReleaseWorkflow = true
		features.GitHub.UseDependabot = true
		features.GitHub.GenerateChangelog = true
	}
	
	return features
}

// NewConfig returns a new configuration for the given project name and type
func NewConfig(projectName string, projectType ProjectType) *Config {
	features := DefaultFeaturesForType(projectType)
	
	return &Config{
		ProjectName: projectName,
		ProjectType: projectType,
		Features:    features,
	}
}

// NewDefaultConfig returns a new default configuration
func NewDefaultConfig() *Config {
	return NewConfig("", DefaultType)
}

// ValidateConfig validates that the configuration is complete and correct
func ValidateConfig(cfg *Config) error {
	return cfg.Validate()
}
