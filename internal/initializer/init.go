package initializer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Caezarr-OSS/Scotter/internal/generator/changelog"
	"github.com/Caezarr-OSS/Scotter/internal/generator/code"
	"github.com/Caezarr-OSS/Scotter/internal/generator/github"
	"github.com/Caezarr-OSS/Scotter/internal/generator/structure"
	"github.com/Caezarr-OSS/Scotter/internal/generator/taskfile"
	"github.com/Caezarr-OSS/Scotter/internal/model"
	"github.com/Caezarr-OSS/Scotter/internal/prompt"
)

// InitProject initializes a new Go project with the specified configuration
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

	// Create the generators
	structureGen := structure.NewGenerator(cfg)
	githubGen := github.NewGenerator(cfg, templatesDir)
	taskfileGen := taskfile.NewGenerator(cfg, templatesDir)
	codeGen := code.NewGenerator(cfg, templatesDir)
	changelogGen := changelog.NewGenerator(cfg)

	// Generate the project structure
	if err := structureGen.Generate(); err != nil {
		return fmt.Errorf("failed to generate project structure: %w", err)
	}

	// Generate .gitignore
	if err := structureGen.GenerateGitIgnore(); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}

	// Generate GitHub configuration
	if err := githubGen.Generate(); err != nil {
		return fmt.Errorf("failed to generate GitHub configuration: %w", err)
	}

	// Generate Taskfile
	if err := taskfileGen.Generate(); err != nil {
		return fmt.Errorf("failed to generate Taskfile: %w", err)
	}

	// Generate CHANGELOG and commitlint configuration
	if err := changelogGen.Generate(); err != nil {
		return fmt.Errorf("failed to generate changelog files: %w", err)
	}

	// Generate Git commit-msg hook for commitlint if enabled
	if cfg.Features.GitHub.UseCommitLint {
		if err := changelogGen.GenerateCommitMsgHook(); err != nil {
			fmt.Printf("Warning: failed to generate Git hook: %v\n", err)
		}
	}

	// Generate code files
	if err := codeGen.Generate(); err != nil {
		return fmt.Errorf("failed to generate code files: %w", err)
	}

	// Success message
	fmt.Println("\n✓ Project initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Review the generated files")
	fmt.Println("2. Run 'go mod tidy' to update dependencies")
	if cfg.Features.UseTaskFile {
		fmt.Println("3. Run 'task build' to build the project")
	} else {
		fmt.Println("3. Run 'go build' to build the project")
	}
	
	return nil
}
