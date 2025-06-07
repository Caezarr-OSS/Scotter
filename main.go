package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Caezarr-OSS/Scotter/internal/config"
	"github.com/Caezarr-OSS/Scotter/internal/generator/taskfile"
	"github.com/Caezarr-OSS/Scotter/internal/initializer"
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// Version information (will be set during build)
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

func main() {
	fmt.Printf("Scotter v%s (%s) built on %s\n", Version, CommitSHA, BuildDate)
	fmt.Println("A modular project bootstrapper with pipeline features")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "init":
		// Initialize project using command line arguments
		fmt.Println("Initializing project...")
		args := config.ParseInitFlags(os.Args[2:])
		
		var cfg *model.Config
		if args.NoInteractive || args.ProjectName != "" || args.GoProjectType != "" || len(args.PipelineFeatures) > 0 {
			// Create configuration from command line args in non-interactive mode
			cfg = config.CreateConfigFromArgs(args)
		} else {
			// Fall back to interactive mode if no specific flags are provided
			fmt.Println("Starting interactive mode. Use --no-interactive flag to disable.")
			if err := initializer.InitProject(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
		
		// Initialize using the configuration
		if err := initializer.InitProjectWithConfig(cfg); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Project initialized successfully!")

	case "target":
		// Need at least one subcommand
		if len(os.Args) < 3 {
			fmt.Println("Error: target command requires a subcommand (add, remove, list)")
			printTargetUsage()
			os.Exit(1)
		}
		
		subCmd := os.Args[2]
		switch subCmd {
		case "add":
			args := config.ParseTargetFlags(subCmd, os.Args[3:])
			if args.OS == "" || args.Arch == "" {
				fmt.Println("Error: --os and --arch are required for target add")
				printTargetUsage()
				os.Exit(1)
			}
			
			// Create and validate build target
			target := model.BuildTarget{OS: args.OS, Arch: args.Arch}
			if !model.ValidOS(target.OS) || !model.ValidArch(target.Arch) {
				fmt.Printf("Error: Invalid OS or architecture: %s/%s\n", target.OS, target.Arch)
				fmt.Println("Valid OS: linux, darwin, windows")
				fmt.Println("Valid architectures: amd64, arm64, 386")
				os.Exit(1)
			}
			
			// Add build target to Taskfile
			targetMgr := taskfile.NewBuildTargetManager(args.ProjectPath)
			err := targetMgr.AddBuildTarget(target)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Build target added: %s/%s\n", target.OS, target.Arch)
			
		case "remove":
			args := config.ParseTargetFlags(subCmd, os.Args[3:])
			if args.OS == "" || args.Arch == "" {
				fmt.Println("Error: --os and --arch are required for target remove")
				printTargetUsage()
				os.Exit(1)
			}
			
			// Create and validate build target
			target := model.BuildTarget{OS: args.OS, Arch: args.Arch}
			
			// Remove build target from Taskfile
			targetMgr := taskfile.NewBuildTargetManager(args.ProjectPath)
			err := targetMgr.RemoveBuildTarget(target)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Build target removed: %s/%s\n", target.OS, target.Arch)
			
		case "list":
			args := config.ParseTargetFlags(subCmd, os.Args[3:])
			
			// List build targets in Taskfile
			targetMgr := taskfile.NewBuildTargetManager(args.ProjectPath)
			targets, err := targetMgr.ListBuildTargets()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println("Build targets:")
			if len(targets) == 0 {
				fmt.Println("  No build targets defined")
			} else {
				for _, target := range targets {
					fmt.Printf("  %s/%s\n", target.OS, target.Arch)
				}
			}
			
		default:
			fmt.Printf("Unknown target subcommand: %s\n", subCmd)
			printTargetUsage()
			os.Exit(1)
		}
		
	case "version":
		fmt.Printf("Scotter v%s (%s) built on %s\n", Version, CommitSHA, BuildDate)
		
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: scotter <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  init           Initialize a project with customizable pipeline features")
	fmt.Println("  target         Manage build targets for Go projects")
	fmt.Println("  version        Show version information")
	fmt.Println("\nOptions for 'init':")
	fmt.Println("  --name           Project name")
	fmt.Println("  --lang           Programming language (go, none) [default: go]")
	fmt.Println("  --go-type        Go project type (default, library, cli, api) [default: default]")
	fmt.Println("  --module         Go module path (e.g., github.com/username/project)")
	fmt.Println("  --features       Comma-separated pipeline features (ci,commit-lint,changelog,release,dependabot,container)")
	fmt.Println("  --github-actions Enable GitHub Actions [default: true]")
	fmt.Println("  --taskfile       Include a Taskfile [default: true]")
	fmt.Println("  --makefile       Include a Makefile [default: false]")
	fmt.Println("  --no-interactive Disable interactive prompts")
	fmt.Println("  --os             Comma-separated list of target operating systems [default: linux,darwin]")
	fmt.Println("  --arch           Comma-separated list of target architectures [default: amd64]")
	fmt.Println("\nRun 'scotter target --help' for information on target commands")
	fmt.Println("\nSupported languages:")
	fmt.Println("  - Go")
	fmt.Println("  - Shell/Script (No specific language)")
	fmt.Println("\nPipeline features:")
	fmt.Println("  - ci           CI Pipeline")
	fmt.Println("  - commit-lint  Enforce conventional commit format")
	fmt.Println("  - changelog    Automated changelog generation")
	fmt.Println("  - release      Automated release pipeline")
	fmt.Println("  - dependabot   Dependency update automation")
	fmt.Println("  - container    Container build configuration")
}

func printTargetUsage() {
	fmt.Println("Usage: scotter target <subcommand> [options]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  add            Add a build target")
	fmt.Println("  remove         Remove a build target")
	fmt.Println("  list           List all build targets")
	fmt.Println("\nOptions:")
	fmt.Println("  --project       Project path [default: current directory]")
	fmt.Println("  --os            Target operating system (linux, darwin, windows)")
	fmt.Println("  --arch          Target architecture (amd64, arm64, 386)")
	fmt.Println("\nExamples:")
	fmt.Println("  scotter target add --os linux --arch arm64")
	fmt.Println("  scotter target remove --os windows --arch amd64")
	fmt.Println("  scotter target list")
}
