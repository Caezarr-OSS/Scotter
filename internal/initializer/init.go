package initializer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Caezarr-OSS/Scotter/internal/generator/changelog"
	"github.com/Caezarr-OSS/Scotter/internal/generator/code"
	"github.com/Caezarr-OSS/Scotter/internal/generator/container"
	"github.com/Caezarr-OSS/Scotter/internal/generator/github"
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

// InitProject initializes a new project with the specified configuration
func InitProject() error {
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

	// Create a new prompter to ask for project configuration
	prompter := prompt.NewProjectPrompt()
	
	// Get project configuration from user input
	cfg := prompter.CollectConfig()
	
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

	// Step 2: Generate pipeline features if GitHub Actions is enabled
	if cfg.Pipeline.UseGitHubActions {
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
	// Create feature generators map
	featureGenerators := map[string]func(*model.Config, string) error{
		"commit-lint": generateCommitLint,
		"changelog": generateChangelog,
		"release": generateRelease,
		"dependabot": generateDependabot,
		"ci": generateCI,
		"container": generateContainer,
	}

	// Create GitHub workflows directory
	if err := os.MkdirAll(".github/workflows", 0755); err != nil {
		return fmt.Errorf("failed to create GitHub workflows directory: %w", err)
	}

	// Generate each selected feature
	for _, featureID := range cfg.Pipeline.SelectedFeatures {
		generatorFunc, exists := featureGenerators[featureID]
		if !exists {
			fmt.Printf("Warning: Unknown feature '%s' selected\n", featureID)
			continue
		}

		if err := generatorFunc(cfg, templatesDir); err != nil {
			return fmt.Errorf("failed to generate feature '%s': %w", featureID, err)
		}
	}

	return nil
}

// generateCommitLint generates commit-lint configuration
func generateCommitLint(cfg *model.Config, templatesDir string) error {
	changelogGen := changelog.NewGenerator(cfg)

	// Generate commitlint configuration
	if err := changelogGen.GenerateCommitLintConfig(); err != nil {
		return fmt.Errorf("failed to generate commitlint config: %w", err)
	}

	// Generate Git commit-msg hook
	if err := changelogGen.GenerateCommitMsgHook(); err != nil {
		fmt.Printf("Warning: failed to generate Git hook: %v\n", err)
	}

	// Generate GitHub workflow for commit validation
	githubGen := github.NewGenerator(cfg, templatesDir)
	if err := githubGen.GenerateCommitLintWorkflow(); err != nil {
		return fmt.Errorf("failed to generate commit-lint workflow: %w", err)
	}

	return nil
}

// generateChangelog generates changelog configuration
func generateChangelog(cfg *model.Config, templatesDir string) error {
	changelogGen := changelog.NewGenerator(cfg)

	// Generate CHANGELOG.md
	if err := changelogGen.GenerateChangelog(); err != nil {
		return fmt.Errorf("failed to generate changelog: %w", err)
	}

	// Generate GitHub workflow for changelog updates
	githubGen := github.NewGenerator(cfg, templatesDir)
	if err := githubGen.GenerateChangelogWorkflow(); err != nil {
		return fmt.Errorf("failed to generate changelog workflow: %w", err)
	}

	return nil
}

// generateRelease generates release configuration
func generateRelease(cfg *model.Config, templatesDir string) error {
	// Generate GitHub workflow for releases
	githubGen := github.NewGenerator(cfg, templatesDir)
	if err := githubGen.GenerateReleaseWorkflow(); err != nil {
		return fmt.Errorf("failed to generate release workflow: %w", err)
	}

	return nil
}

// generateDependabot generates Dependabot configuration
func generateDependabot(cfg *model.Config, templatesDir string) error {
	// Generate Dependabot configuration
	githubGen := github.NewGenerator(cfg, templatesDir)
	if err := githubGen.GenerateDependabotConfig(); err != nil {
		return fmt.Errorf("failed to generate Dependabot config: %w", err)
	}

	return nil
}

// generateCI generates CI pipeline configuration
func generateCI(cfg *model.Config, templatesDir string) error {
	// Create GitHub generator
	githubGen := github.NewGenerator(cfg, templatesDir)

	// Generate CI workflow
	if err := githubGen.GenerateCIWorkflow(); err != nil {
		return err
	}

	return nil
}

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
