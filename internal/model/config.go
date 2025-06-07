package model

import (
	"fmt"
)

// LanguageType defines the programming language for the project
type LanguageType string

const (
	// GoLang is for Go projects
	GoLang LanguageType = "go"
	// NoLang is for projects without a specific language (shell scripts, etc.)
	NoLang LanguageType = "none"
)

// GoProjectType defines the type of Go project to create
type GoProjectType string

const (
	// DefaultGoType is a minimal Go project structure
	DefaultGoType GoProjectType = "default"
	// LibraryGoType is for Go packages/libraries
	LibraryGoType GoProjectType = "library"
	// CLIGoType is for command-line applications
	CLIGoType GoProjectType = "cli"
	// APIGoType is for HTTP API/service applications
	APIGoType GoProjectType = "api"
	// CompleteGoType is for a complete project with all features
	CompleteGoType GoProjectType = "complete"
)

// String returns the string representation of GoProjectType
func (pt GoProjectType) String() string {
	switch pt {
	case DefaultGoType:
		return "Default"
	case LibraryGoType:
		return "Library"
	case CLIGoType:
		return "CLI"
	case APIGoType:
		return "API"
	case CompleteGoType:
		return "Complete"
	default:
		return "Unknown"
	}
}

// PipelineFeature represents a feature that can be added to a pipeline
type PipelineFeature struct {
	// ID is the unique identifier for the feature
	ID string
	// Name is the display name of the feature
	Name string
	// Description explains what the feature does
	Description string
	// Dependencies lists other features that this feature requires
	Dependencies []string
}

// GoConfig holds Go-specific configuration
type GoConfig struct {
	// ModulePath is the Go module path (e.g., "github.com/username/project")
	ModulePath string
	// ProjectType determines the structure and features
	ProjectType GoProjectType
	// UseTaskFile specifies whether to include a Taskfile
	UseTaskFile bool
	// UseMakeFile specifies whether to include a Makefile
	UseMakeFile bool
	// BuildTargets specifies the OS/architecture combinations to target
	BuildTargets []BuildTarget
}

// ContainerFileFormat defines the format of container configuration files
type ContainerFileFormat string

const (
	// DockerfileFormat is the standard Docker format
	DockerfileFormat ContainerFileFormat = "dockerfile"
	// ContainerfileFormat is the OCI/Podman format
	ContainerfileFormat ContainerFileFormat = "containerfile"
)

// String returns the string representation of ContainerFileFormat
func (f ContainerFileFormat) String() string {
	switch f {
	case DockerfileFormat:
		return "Dockerfile"
	case ContainerfileFormat:
		return "Containerfile"
	default:
		return "Dockerfile"
	}
}

// PipelineConfig holds pipeline configuration
type PipelineConfig struct {
	// UseGitHubActions specifies whether to include GitHub Actions workflows
	UseGitHubActions bool
	// SelectedFeatures contains the IDs of selected pipeline features
	SelectedFeatures []string
	// ContainerFormat specifies the format for container configuration files
	ContainerFormat ContainerFileFormat
}

// Config holds the project configuration
type Config struct {
	// ProjectName is the name of the project
	ProjectName string
	// Language is the programming language for the project
	Language LanguageType
	// Go contains Go-specific configuration (used when Language == GoLang)
	Go GoConfig
	// Pipeline contains pipeline configuration
	Pipeline PipelineConfig
	// Directories to create
	Directories []string
}

// Validate validates that the configuration is complete and correct
func (cfg *Config) Validate() error {
	// Check required fields
	if cfg.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Validate language type
	validLanguage := false
	for _, lang := range []LanguageType{GoLang, NoLang} {
		if cfg.Language == lang {
			validLanguage = true
			break
		}
	}

	if !validLanguage {
		return fmt.Errorf("invalid language type: %s", cfg.Language)
	}

	// Validate language-specific configuration
	if cfg.Language == GoLang {
		// Check project type
		validType := false
		for _, pt := range []GoProjectType{DefaultGoType, LibraryGoType, CLIGoType, APIGoType} {
			if cfg.Go.ProjectType == pt {
				validType = true
				break
			}
		}

		if !validType {
			return fmt.Errorf("invalid Go project type: %s", cfg.Go.ProjectType)
		}
	}

	return nil
}

// NewConfig returns a new configuration with default values
func NewConfig() *Config {
	return &Config{
		ProjectName: "",
		Language:    GoLang,
		Go: GoConfig{
			ProjectType: DefaultGoType,
			UseTaskFile: true,
			BuildTargets: []BuildTarget{
				{OS: "linux", Arch: "amd64"},
				{OS: "darwin", Arch: "amd64"},
			},
		},
		Pipeline: PipelineConfig{
			UseGitHubActions:  true,
			SelectedFeatures: []string{"ci"},
			ContainerFormat:  DockerfileFormat,
		},
		Directories: []string{},
	}
}

// ValidateConfig validates that the configuration is complete and correct
func ValidateConfig(cfg *Config) error {
	return cfg.Validate()
}

// AvailablePipelineFeatures returns the list of available pipeline features
func AvailablePipelineFeatures() []PipelineFeature {
	return []PipelineFeature{
		{
			ID:          "ci",
			Name:        "CI Pipeline",
			Description: "Continuous Integration pipeline (build, test)",
			Dependencies: []string{},
		},
		{
			ID:          "commit-lint",
			Name:        "Commit Lint",
			Description: "Validates commit message format",
			Dependencies: []string{},
		},
		{
			ID:          "changelog",
			Name:        "Changelog",
			Description: "Generates a changelog from commits",
			Dependencies: []string{"commit-lint"},
		},
		{
			ID:          "release",
			Name:        "Automatic Release",
			Description: "Creates GitHub releases automatically",
			Dependencies: []string{"changelog"},
		},
		{
			ID:          "dependabot",
			Name:        "Dependabot",
			Description: "Automatic dependency updates",
			Dependencies: []string{},
		},
		{
			ID:          "container",
			Name:        "Container",
			Description: "Container support",
			Dependencies: []string{},
		},
	}
}

// ResolveFeatureDependencies ensures all dependencies are included
func ResolveFeatureDependencies(selectedFeatureIDs []string) []string {
	// Create a map for quick lookup
	selected := make(map[string]bool)
	for _, id := range selectedFeatureIDs {
		selected[id] = true
	}

	// Get all available features
	allFeatures := AvailablePipelineFeatures()

	// Resolve dependencies
	changed := true
	for changed {
		changed = false
		
		for _, feature := range allFeatures {
			// If this feature is selected
			if selected[feature.ID] {
				// Check its dependencies
				for _, depID := range feature.Dependencies {
					if !selected[depID] {
						selected[depID] = true
						changed = true
					}
				}
			}
		}
	}

	// Convert back to slice
	var result []string
	for id := range selected {
		result = append(result, id)
	}

	return result
}
