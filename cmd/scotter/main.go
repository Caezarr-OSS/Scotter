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
	fmt.Println("A Go project bootstrapper for GitHub")

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
	fmt.Println("  init      Initialize a Go project with GitHub support")
	fmt.Println("  version   Show version information")
}
