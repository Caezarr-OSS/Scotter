package initializer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Caezarr-OSS/Scotter/internal/generator/ci"
	"github.com/Caezarr-OSS/Scotter/internal/generator/code"
	"github.com/Caezarr-OSS/Scotter/internal/generator/container"
	"github.com/Caezarr-OSS/Scotter/internal/generator/structure"
	"github.com/Caezarr-OSS/Scotter/internal/generator/taskfile"
	"github.com/Caezarr-OSS/Scotter/internal/model"
	"github.com/Caezarr-OSS/Scotter/internal/prompt"
)

// Generator defines the interface for all project generators
type Generator interface {
	Generate() error
}

// PipelineFeatureGenerator generates a specific pipeline feature
type PipelineFeatureGenerator struct {
	ID          string
	Config      *model.Config
	TemplatesDir string
	Generate    func(cfg *model.Config, templatesDir string) error
}

// InitProject initializes a new project with interactive configuration
func InitProject() error {
	// Create a new prompter to ask for project configuration
	prompter := prompt.NewProjectPrompt()
	
	// Get project configuration from user input
	cfg := prompter.CollectConfig()

	// Initialize the project with the collected config
	return InitProjectWithConfig(cfg)
}

// InitProjectWithConfig initializes a new project with the provided configuration
func InitProjectWithConfig(cfg *model.Config) error {
	// Get the executable path to find templates
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Templates are expected to be in the same directory as the executable
	templatesDir := filepath.Join(filepath.Dir(execPath), "templates")
	
	// For development, use a relative path if templates are not found
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		// Try to find templates in the project directory
		templatesDir = filepath.Join("internal", "templates")
		if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
			return fmt.Errorf("templates directory not found: %s", templatesDir)
		}
	}
	
	// Detect OS and set appropriate line endings
	if runtime.GOOS == "windows" {
		fmt.Println("Detected Windows OS, will use CRLF line endings for generated files")
	} else {
		fmt.Println("Detected Unix-like OS, will use LF line endings for generated files")
	}
	
	// Validate the configuration
	if err := model.ValidateConfig(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Step 1: Generate project structure based on language
	if err := generateProjectStructure(cfg, templatesDir); err != nil {
		return err
	}

	// Step 2: Generate pipeline features if a CI system is configured
	if cfg.Pipeline.CIType != "" || cfg.Pipeline.UseGitHubActions {
		if err := generatePipelineFeatures(cfg, templatesDir); err != nil {
			return err
		}
	}

	// Success message
	fmt.Println("\n✓ Project initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Review the generated files")

	// Language-specific next steps
	if cfg.Language == model.GoLang {
		fmt.Println("2. Run 'go mod tidy' to update dependencies")
		if cfg.Go.UseTaskFile {
			fmt.Println("3. Run 'task build' to build the project")
		} else {
			fmt.Println("3. Run 'go build' to build the project")
		}
	}
	
	return nil
}

// generateProjectStructure creates the basic project structure based on language
func generateProjectStructure(cfg *model.Config, templatesDir string) error {
	// Create the structure generator
	structureGen := structure.NewGenerator(cfg)

	// Generate the project structure
	if err := structureGen.Generate(); err != nil {
		return fmt.Errorf("failed to generate project structure: %w", err)
	}

	// Generate .gitignore
	if err := structureGen.GenerateGitIgnore(); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}

	// Language-specific generation
	if cfg.Language == model.GoLang {
		// For Go projects, generate code and build files
		codeGen := code.NewGenerator(cfg, templatesDir)
		if err := codeGen.Generate(); err != nil {
			return fmt.Errorf("failed to generate code files: %w", err)
		}

		// Generate Taskfile if enabled
		if cfg.Go.UseTaskFile {
			taskfileGen := taskfile.NewGenerator(cfg, templatesDir)
			if err := taskfileGen.Generate(); err != nil {
				return fmt.Errorf("failed to generate Taskfile: %w", err)
			}
		}
	}

	return nil
}

// generatePipelineFeatures generates the selected pipeline features
func generatePipelineFeatures(cfg *model.Config, templatesDir string) error {
	// Create the CI Manager factory
	ciFactory := ci.NewCIManagerFactory(templatesDir)
	
	// Create the appropriate CI manager based on the config
	ciManager, err := ciFactory.CreateManager(cfg)
	if err != nil {
		return fmt.Errorf("failed to create CI manager: %w", err)
	}
	
	// If CI is disabled or no manager available, nothing to do
	if ciManager == nil {
		fmt.Println("No CI system configured, skipping pipeline generation")
		return nil
	}
	
	// Create directories based on CI type
	switch ciManager.GetType() {
	case model.GithubActionsCI:
		// Create GitHub workflows directory
		if err := os.MkdirAll(".github/workflows", 0755); err != nil {
			return fmt.Errorf("failed to create GitHub workflows directory: %w", err)
		}
	case model.GitlabCI:
		// GitLab CI doesn't need special directories
	case model.CircleCI:
		// Create CircleCI directory
		if err := os.MkdirAll(".circleci", 0755); err != nil {
			return fmt.Errorf("failed to create CircleCI directory: %w", err)
		}
	}

	// Generate all selected CI features using the CI manager
	if err := ciManager.Generate(); err != nil {
		return fmt.Errorf("failed to generate CI configuration: %w", err)
	}
	
	// Generate container if selected
	if hasFeature(cfg.Pipeline.SelectedFeatures, "container") {
		if err := generateContainer(cfg, templatesDir); err != nil {
			return fmt.Errorf("failed to generate container feature: %w", err)
		}
	}

	return nil
}

// hasFeature checks if a feature is in the selected features list
func hasFeature(features []string, target string) bool {
	for _, f := range features {
		if f == target {
			return true
		}
	}
	return false
}

// Note: Previous functions for generating CI components have been removed
// These are now handled by specialized CI Managers

// generateContainer generates container configuration
func generateContainer(cfg *model.Config, templatesDir string) error {
	// Create container generator
	containerGen := container.NewGenerator(cfg, templatesDir)

	// Generate container files
	if err := containerGen.Generate(); err != nil {
		return err
	}

	return nil
}
