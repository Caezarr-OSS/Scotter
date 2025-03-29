package model

import (
	"testing"
)

// TestProjectTypeString tests the String method of ProjectType
func TestProjectTypeString(t *testing.T) {
	tests := []struct {
		name         string
		projectType  ProjectType
		expectedName string
	}{
		{
			name:         "default type",
			projectType:  DefaultType,
			expectedName: "Default",
		},
		{
			name:         "library type",
			projectType:  LibraryType,
			expectedName: "Library",
		},
		{
			name:         "CLI type",
			projectType:  CLIType,
			expectedName: "CLI",
		},
		{
			name:         "API type",
			projectType:  APIType,
			expectedName: "API",
		},
		{
			name:         "complete type",
			projectType:  CompleteType,
			expectedName: "Complete",
		},
		{
			name:         "unknown type",
			projectType:  ProjectType("unknown"),
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
				ProjectType: DefaultType,
			},
			expectError: false,
		},
		{
			name: "empty project name",
			config: Config{
				ProjectName: "",
				ProjectType: DefaultType,
			},
			expectError: true,
		},
		{
			name: "invalid project type",
			config: Config{
				ProjectName: "testproject",
				ProjectType: ProjectType("invalid"),
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

// TestDefaultFeatures tests the default features for each project type
func TestDefaultFeatures(t *testing.T) {
	tests := []struct {
		name        string
		projectType ProjectType
		checkFunc   func(*testing.T, Features)
	}{
		{
			name:        "default type features",
			projectType: DefaultType,
			checkFunc: func(t *testing.T, f Features) {
				if f.UseTaskFile {
					t.Error("expected UseTaskFile to be false for DefaultType")
				}
				if f.GitHub.UseWorkflows {
					t.Error("expected UseWorkflows to be false for DefaultType")
				}
			},
		},
		{
			name:        "library type features",
			projectType: LibraryType,
			checkFunc: func(t *testing.T, f Features) {
				if !f.GitHub.UseWorkflows {
					t.Error("expected UseWorkflows to be true for LibraryType")
				}
				if !f.GitHub.GenerateChangelog {
					t.Error("expected GenerateChangelog to be true for LibraryType")
				}
			},
		},
		{
			name:        "CLI type features",
			projectType: CLIType,
			checkFunc: func(t *testing.T, f Features) {
				if !f.UseTaskFile {
					t.Error("expected UseTaskFile to be true for CLIType")
				}
				if !f.GitHub.UseWorkflows {
					t.Error("expected UseWorkflows to be true for CLIType")
				}
				if !f.GitHub.UseReleaseWorkflow {
					t.Error("expected UseReleaseWorkflow to be true for CLIType")
				}
			},
		},
		{
			name:        "API type features",
			projectType: APIType,
			checkFunc: func(t *testing.T, f Features) {
				if !f.UseTaskFile {
					t.Error("expected UseTaskFile to be true for APIType")
				}
				if !f.GitHub.UseWorkflows {
					t.Error("expected UseWorkflows to be true for APIType")
				}
			},
		},
		{
			name:        "complete type features",
			projectType: CompleteType,
			checkFunc: func(t *testing.T, f Features) {
				if !f.UseTaskFile {
					t.Error("expected UseTaskFile to be true for CompleteType")
				}
				if !f.GitHub.UseWorkflows {
					t.Error("expected UseWorkflows to be true for CompleteType")
				}
				if !f.GitHub.UseCommitLint {
					t.Error("expected UseCommitLint to be true for CompleteType")
				}
				if !f.GitHub.UseReleaseWorkflow {
					t.Error("expected UseReleaseWorkflow to be true for CompleteType")
				}
				if !f.GitHub.UseDependabot {
					t.Error("expected UseDependabot to be true for CompleteType")
				}
				if !f.GitHub.GenerateChangelog {
					t.Error("expected GenerateChangelog to be true for CompleteType")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			features := DefaultFeaturesForType(tt.projectType)
			tt.checkFunc(t, features)
		})
	}
}

// TestNewConfig tests the NewConfig constructor function
func TestNewConfig(t *testing.T) {
	config := NewConfig("testproject", CLIType)
	
	if config.ProjectName != "testproject" {
		t.Errorf("expected ProjectName to be 'testproject', got %q", config.ProjectName)
	}
	
	if config.ProjectType != CLIType {
		t.Errorf("expected ProjectType to be CLIType, got %v", config.ProjectType)
	}
	
	// Check that it has default features for CLI type
	if !config.Features.UseTaskFile {
		t.Error("expected UseTaskFile to be true for CLI type")
	}
	
	if !config.Features.GitHub.UseWorkflows {
		t.Error("expected UseWorkflows to be true for CLI type")
	}
	
	if !config.Features.GitHub.UseReleaseWorkflow {
		t.Error("expected UseReleaseWorkflow to be true for CLI type")
	}
}
