package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// ProjectPrompt manages interactive user prompts for project configuration
type ProjectPrompt struct {
	reader *bufio.Reader
}

// NewProjectPrompt creates a new project prompt manager
func NewProjectPrompt() *ProjectPrompt {
	return &ProjectPrompt{
		reader: bufio.NewReader(os.Stdin),
	}
}

// AskString asks the user for a string input
func (p *ProjectPrompt) AskString(question string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", question, defaultValue)
	} else {
		fmt.Printf("%s: ", question)
	}

	answer, _ := p.reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	if answer == "" {
		return defaultValue
	}
	return answer
}

// AskBool asks the user for a yes/no answer
func (p *ProjectPrompt) AskBool(question string, defaultValue bool) bool {
	defaultStr := "n"
	if defaultValue {
		defaultStr = "y"
	}

	for {
		answer := p.AskString(fmt.Sprintf("%s (y/n)", question), defaultStr)
		answer = strings.ToLower(answer)

		if answer == "y" || answer == "yes" {
			return true
		} else if answer == "n" || answer == "no" {
			return false
		} else if answer == "" {
			return defaultValue
		}

		fmt.Println("Please answer 'y' or 'n'")
	}
}

// AskProjectType asks the user for the project type
func (p *ProjectPrompt) AskProjectType() model.ProjectType {
	fmt.Println("\nProject types:")
	fmt.Println("1. Default - A minimal structure for simple projects or scripting")
	fmt.Println("2. Library - For reusable Go packages")
	fmt.Println("3. CLI - For command-line applications")
	fmt.Println("4. API - For HTTP API/service applications")
	fmt.Println("5. Complete - All features enabled")

	for {
		answer := p.AskString("Select a project type (1-5)", "1")

		switch answer {
		case "1":
			return model.DefaultType
		case "2":
			return model.LibraryType
		case "3":
			return model.CLIType
		case "4":
			return model.APIType
		case "5":
			return model.CompleteType
		default:
			fmt.Println("Please enter a number between 1 and 5")
		}
	}
}

// CollectConfig prompts the user for project configuration
func (p *ProjectPrompt) CollectConfig() *model.Config {
	cfg := model.NewDefaultConfig()

	// Basic project info
	fmt.Println("\n=== Basic Project Configuration ===")
	cfg.ProjectName = p.AskString("Project name", cfg.ProjectName)
	cfg.ModulePath = p.AskString("Go module path (e.g., github.com/username/project)", fmt.Sprintf("github.com/username/%s", cfg.ProjectName))
	cfg.ProjectType = p.AskProjectType()

	// GitHub configuration
	fmt.Println("\n=== GitHub Configuration ===")
	cfg.Features.GitHub.UseWorkflows = p.AskBool("Include GitHub Actions workflows", true)
	
	if cfg.Features.GitHub.UseWorkflows {
		cfg.Features.GitHub.UseCommitLint = p.AskBool("Include commit message validation", true)
		cfg.Features.GitHub.UseReleaseWorkflow = p.AskBool("Include automatic release workflow", true)
		cfg.Features.GitHub.UseDependabot = p.AskBool("Include Dependabot configuration", true)
	}

	// Build tools
	fmt.Println("\n=== Build Tools ===")
	cfg.Features.UseTaskFile = p.AskBool("Include Taskfile", true)
	cfg.Features.UseMakeFile = p.AskBool("Include Makefile", false)

	return cfg
}
