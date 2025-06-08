package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// CommandLineArgs represents parsed command line arguments
type CommandLineArgs struct {
	// Common flags
	ProjectPath string

	// Init command flags
	ProjectName       string
	Language          string
	GoProjectType     string
	GoModulePath      string
	PipelineFeatures  []string
	UseGitHubActions  bool
	UseTaskfile       bool
	UseMakefile       bool
	NoInteractive     bool
	ContainerFormat   string
	TargetOS          []string
	TargetArch        []string

	// Target command flags
	TargetAction string
	OS           string
	Arch         string
}

// ParseInitFlags parses command-line arguments for the init command
func ParseInitFlags(args []string) *CommandLineArgs {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	
	// Define flags
	projectName := initCmd.String("name", "", "Project name")
	language := initCmd.String("lang", "go", "Programming language (go, none)")
	goType := initCmd.String("go-type", "default", "Go project type (default, local-library, distributed-library, cli, api)")
	modulePath := initCmd.String("module", "", "Go module path (e.g., github.com/username/project)")
	features := initCmd.String("features", "ci", "Comma-separated pipeline features (ci,commit-lint,changelog,release,dependabot,container)")
	githubActions := initCmd.Bool("github-actions", true, "Configure GitHub Actions")
	useTaskfile := initCmd.Bool("taskfile", true, "Include a Taskfile")
	useMakefile := initCmd.Bool("makefile", false, "Include a Makefile")
	// Interactive mode has been completely removed - always non-interactive
	containerFormat := initCmd.String("container-format", "dockerfile", "Container file format (dockerfile, containerfile)")
	targetOS := initCmd.String("os", "linux,darwin", "Comma-separated list of target operating systems")
	targetArch := initCmd.String("arch", "amd64", "Comma-separated list of target architectures")
	
	// Parse flags
	err := initCmd.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		// Let the caller handle the exit
		return nil
	}

	// Split comma-separated values
	var featuresList []string
	if *features != "" {
		featuresList = strings.Split(*features, ",")
	}

	var osList []string
	if *targetOS != "" {
		osList = strings.Split(*targetOS, ",")
	}

	var archList []string
	if *targetArch != "" {
		archList = strings.Split(*targetArch, ",")
	}
	
	return &CommandLineArgs{
		ProjectName:      *projectName,
		Language:         *language,
		GoProjectType:    *goType,
		GoModulePath:     *modulePath,
		PipelineFeatures: featuresList,
		UseGitHubActions: *githubActions,
		UseTaskfile:      *useTaskfile,
		UseMakefile:      *useMakefile,
		NoInteractive:    true, // Mode non-interactif toujours actif
		ContainerFormat:  *containerFormat,
		TargetOS:         osList,
		TargetArch:       archList,
	}
}

// ParseTargetFlags parses command-line arguments for the target commands
func ParseTargetFlags(command string, args []string) *CommandLineArgs {
	targetCmd := flag.NewFlagSet("target", flag.ExitOnError)
	
	// Define flags
	projectPath := targetCmd.String("project", ".", "Project path")
	os := targetCmd.String("os", "", "Target operating system (linux, darwin, windows)")
	arch := targetCmd.String("arch", "", "Target architecture (amd64, arm64, 386)")
	
	// Parse flags
	err := targetCmd.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		// Let the caller handle the exit
		return nil
	}
	
	return &CommandLineArgs{
		ProjectPath:  *projectPath,
		TargetAction: command,
		OS:           *os,
		Arch:         *arch,
	}
}

// CreateConfigFromArgs creates a Config object from command-line arguments
func CreateConfigFromArgs(args *CommandLineArgs) *model.Config {
	cfg := model.NewConfig()
	
	// Set basic properties
	if args.ProjectName != "" {
		cfg.ProjectName = args.ProjectName
	}
	
	// Set language
	switch args.Language {
	case "go":
		cfg.Language = model.GoLang
	case "none":
		cfg.Language = model.NoLang
	default:
		cfg.Language = model.GoLang
	}
	
	// Set Go-specific config
	if cfg.Language == model.GoLang {
		// Set project type
		switch args.GoProjectType {
		case "default":
			cfg.Go.ProjectType = model.DefaultGoType
		case "library":
			// Pour maintenir la compatibilité avec les versions antérieures
			cfg.Go.ProjectType = model.LibraryGoType
		case "local-library":
			cfg.Go.ProjectType = model.LocalLibraryGoType
		case "distributed-library":
			cfg.Go.ProjectType = model.DistributedLibraryGoType
		case "cli":
			cfg.Go.ProjectType = model.CLIGoType
		case "api":
			cfg.Go.ProjectType = model.APIGoType
		default:
			cfg.Go.ProjectType = model.DefaultGoType
		}
		
		// Set module path
		if args.GoModulePath != "" {
			cfg.Go.ModulePath = args.GoModulePath
		} else if args.ProjectName != "" {
			cfg.Go.ModulePath = fmt.Sprintf("github.com/username/%s", args.ProjectName)
		}
		
		// Set taskfile/makefile preferences
		cfg.Go.UseTaskFile = args.UseTaskfile
		cfg.Go.UseMakeFile = args.UseMakefile
		
		// Set build targets
		cfg.Go.BuildTargets = []model.BuildTarget{}
		for _, osName := range args.TargetOS {
			for _, archName := range args.TargetArch {
				// Validate OS and architecture
				if model.ValidOS(osName) && model.ValidArch(archName) {
					cfg.Go.BuildTargets = append(cfg.Go.BuildTargets, model.BuildTarget{
						OS:   osName,
						Arch: archName,
					})
				}
			}
		}
		
		// Use defaults if no valid targets specified
		if len(cfg.Go.BuildTargets) == 0 {
			cfg.Go.BuildTargets = []model.BuildTarget{
				{OS: "linux", Arch: "amd64"},
			}
		}
	}
	
	// Set pipeline config
	cfg.Pipeline.UseGitHubActions = args.UseGitHubActions
	
	// Set pipeline features
	if len(args.PipelineFeatures) > 0 {
		cfg.Pipeline.SelectedFeatures = model.ResolveFeatureDependencies(args.PipelineFeatures)
	}
	
	// Set container format if container feature is selected
	if args.ContainerFormat != "" {
		switch args.ContainerFormat {
		case "dockerfile":
			cfg.Pipeline.ContainerFormat = model.DockerfileFormat
		case "containerfile":
			cfg.Pipeline.ContainerFormat = model.ContainerfileFormat
		}
	}
	
	return cfg
}
