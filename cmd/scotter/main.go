package main

import (
	"fmt"
	"os"

	"github.com/Caezarr-OSS/Scotter/internal/initializer"
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
		fmt.Println("Initializing project...")
		if err := initializer.InitProject(); err != nil {
			fmt.Printf("Error: %v\n", err)
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
	fmt.Println("Usage: scotter <command>")
	fmt.Println("Commands:")
	fmt.Println("  init      Initialize a project with customizable pipeline features")
	fmt.Println("  version   Show version information")
	fmt.Println("\nSupported languages:")
	fmt.Println("  - Go")
	fmt.Println("  - Shell/Script (No specific language)")
	fmt.Println("\nPipeline features:")
	fmt.Println("  - Commit Lint")
	fmt.Println("  - Changelog Generation")
	fmt.Println("  - Automatic Release")
	fmt.Println("  - Dependabot")
	fmt.Println("  - CI Pipeline")
}
