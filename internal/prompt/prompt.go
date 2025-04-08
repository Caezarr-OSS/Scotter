package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

// AskSelect asks the user to select an option from a list
func (p *ProjectPrompt) AskSelect(question string, options []string, defaultIndex int) int {
	fmt.Println("\n" + question + ":")
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	defaultStr := ""
	if defaultIndex >= 0 && defaultIndex < len(options) {
		defaultStr = strconv.Itoa(defaultIndex + 1)
	}

	for {
		answer := p.AskString(fmt.Sprintf("Select an option (1-%d)", len(options)), defaultStr)
		if answer == "" && defaultStr != "" {
			return defaultIndex
		}

		index, err := strconv.Atoi(answer)
		if err == nil && index >= 1 && index <= len(options) {
			return index - 1
		}

		fmt.Printf("Please enter a number between 1 and %d\n", len(options))
	}
}

// AskMultiSelect asks the user to select multiple options from a list
func (p *ProjectPrompt) AskMultiSelect(question string, options []model.PipelineFeature) []string {
	// Display options
	fmt.Println("\n" + question + ":")
	for i, option := range options {
		fmt.Printf("%d. %s - %s\n", i+1, option.Name, option.Description)
	}

	// Get user input
	selectedIDs := []string{}
	for {
		answer := p.AskString("Enter numbers separated by commas (e.g., 1,3,4) or 'all' for all features", "all")
		answer = strings.ToLower(answer)

		// Handle 'all' option
		if answer == "all" {
			for _, option := range options {
				selectedIDs = append(selectedIDs, option.ID)
			}
			break
		}

		// Parse comma-separated list
		valid := true
		parts := strings.Split(answer, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}

			index, err := strconv.Atoi(trimmed)
			if err != nil || index < 1 || index > len(options) {
				fmt.Printf("Invalid option: %s. Please enter numbers between 1 and %d\n", trimmed, len(options))
				valid = false
				break
			}

			selectedIDs = append(selectedIDs, options[index-1].ID)
		}

		if valid {
			break
		}
	}

	// Resolve dependencies
	return model.ResolveFeatureDependencies(selectedIDs)
}

// AskLanguage asks the user for the project language
func (p *ProjectPrompt) AskLanguage() model.LanguageType {
	options := []string{
		"Go - For Go projects",
		"None/Shell - For script projects or multi-language projects",
	}

	index := p.AskSelect("Select project language", options, 0)
	if index == 0 {
		return model.GoLang
	}
	return model.NoLang
}

// AskGoProjectType asks the user for the Go project type
func (p *ProjectPrompt) AskGoProjectType() model.GoProjectType {
	options := []string{
		"Default - A minimal structure for simple projects",
		"Library - For reusable Go packages",
		"CLI - For command-line applications",
		"API - For HTTP API/service applications",
	}

	index := p.AskSelect("Select Go project type", options, 0)
	switch index {
	case 0:
		return model.DefaultGoType
	case 1:
		return model.LibraryGoType
	case 2:
		return model.CLIGoType
	case 3:
		return model.APIGoType
	default:
		return model.DefaultGoType
	}
}

// AskContainerFileFormat asks the user for their preferred container file format
func (p *ProjectPrompt) AskContainerFileFormat() model.ContainerFileFormat {
	options := []string{
		"Dockerfile (Docker standard)",
		"Containerfile (Podman/OCI standard)",
	}

	index := p.AskSelect("Select your preferred container file format", options, 0)
	switch index {
	case 0:
		return model.DockerfileFormat
	case 1:
		return model.ContainerfileFormat
	default:
		return model.DockerfileFormat
	}
}

// CollectConfig prompts the user for project configuration
func (p *ProjectPrompt) CollectConfig() *model.Config {
	cfg := model.NewConfig()

	// Basic project info
	fmt.Println("\n=== Basic Project Configuration ===")
	cfg.ProjectName = p.AskString("Project name", cfg.ProjectName)
	
	// Language selection
	cfg.Language = p.AskLanguage()

	// Language-specific configuration
	if cfg.Language == model.GoLang {
		fmt.Println("\n=== Go Configuration ===")
		cfg.Go.ModulePath = p.AskString("Go module path (e.g., github.com/username/project)", 
			fmt.Sprintf("github.com/username/%s", cfg.ProjectName))
		cfg.Go.ProjectType = p.AskGoProjectType()
		cfg.Go.UseTaskFile = p.AskBool("Include Taskfile", true)
		cfg.Go.UseMakeFile = p.AskBool("Include Makefile", false)
	}

	// Pipeline configuration
	fmt.Println("\n=== Pipeline Configuration ===")
	cfg.Pipeline.UseGitHubActions = p.AskBool("Configure GitHub Actions", true)
	
	if cfg.Pipeline.UseGitHubActions {
		// Display available pipeline features and let user select
		fmt.Println("\nAvailable pipeline features:")
		features := model.AvailablePipelineFeatures()
		selectedFeatures := p.AskMultiSelect("Select pipeline features", features)
		
		// Store selected features
		cfg.Pipeline.SelectedFeatures = selectedFeatures
		
		// Show selected features with dependencies resolved
		fmt.Println("\nSelected features (including dependencies):")
		featureMap := make(map[string]model.PipelineFeature)
		for _, f := range features {
			featureMap[f.ID] = f
		}
		
		for _, id := range cfg.Pipeline.SelectedFeatures {
			if feature, ok := featureMap[id]; ok {
				fmt.Printf("- %s\n", feature.Name)
			}
		}
		
		// If container feature is selected, ask for container file format
		for _, id := range cfg.Pipeline.SelectedFeatures {
			if id == "container" {
				fmt.Println("\n=== Container Configuration ===")
				cfg.Pipeline.ContainerFormat = p.AskContainerFileFormat()
				break
			}
		}
	}

	return cfg
}
