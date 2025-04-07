package model

import (
	"testing"
)

// TestGoProjectTypeString tests the String method of GoProjectType
func TestGoProjectTypeString(t *testing.T) {
	tests := []struct {
		name         string
		projectType  GoProjectType
		expectedName string
	}{
		{
			name:         "default type",
			projectType:  DefaultGoType,
			expectedName: "Default",
		},
		{
			name:         "library type",
			projectType:  LibraryGoType,
			expectedName: "Library",
		},
		{
			name:         "CLI type",
			projectType:  CLIGoType,
			expectedName: "CLI",
		},
		{
			name:         "API type",
			projectType:  APIGoType,
			expectedName: "API",
		},
		{
			name:         "unknown type",
			projectType:  GoProjectType("unknown"),
			expectedName: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.projectType.String()
			if result != tt.expectedName {
				t.Errorf("expected %q, got %q", tt.expectedName, result)
			}
		})
	}
}

// TestConfigValidation tests the validation of Config
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				ProjectName: "testproject",
				Language: GoLang,
				Go: GoConfig{
					ProjectType: DefaultGoType,
				},
			},
			expectError: false,
		},
		{
			name: "empty project name",
			config: Config{
				ProjectName: "",
				Language: GoLang,
				Go: GoConfig{
					ProjectType: DefaultGoType,
				},
			},
			expectError: true,
		},
		{
			name: "invalid language",
			config: Config{
				ProjectName: "testproject",
				Language: LanguageType("invalid"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("expected error? %v, got error: %v", tt.expectError, err)
			}
		})
	}
}

// TestPipelineFeatures tests the available pipeline features
func TestPipelineFeatures(t *testing.T) {
	features := AvailablePipelineFeatures()
	
	// Verify that we have the expected number of features
	if len(features) < 3 {
		t.Errorf("expected at least 3 pipeline features, got %d", len(features))
	}
	
	// Check for specific features
	foundCommitLint := false
	foundChangelog := false
	foundCI := false
	
	for _, f := range features {
		switch f.ID {
		case "commit-lint":
			foundCommitLint = true
		case "changelog":
			foundChangelog = true
			// Check that changelog depends on commit-lint
			hasDependency := false
			for _, dep := range f.Dependencies {
				if dep == "commit-lint" {
					hasDependency = true
					break
				}
			}
			if !hasDependency {
				t.Error("expected changelog to depend on commit-lint")
			}
		case "ci":
			foundCI = true
		}
	}
	
	if !foundCommitLint {
		t.Error("commit-lint feature not found")
	}
	if !foundChangelog {
		t.Error("changelog feature not found")
	}
	if !foundCI {
		t.Error("ci feature not found")
	}
}

// TestNewConfig tests the NewConfig constructor function
func TestNewConfig(t *testing.T) {
	config := NewConfig()
	
	// Check that it has default values
	if config.ProjectName != "" {
		t.Errorf("expected ProjectName to be empty, got %q", config.ProjectName)
	}
	
	if config.Language != GoLang {
		t.Errorf("expected Language to be GoLang, got %v", config.Language)
	}
	
	if config.Go.ProjectType != DefaultGoType {
		t.Errorf("expected Go.ProjectType to be DefaultGoType, got %v", config.Go.ProjectType)
	}
	

	

}
