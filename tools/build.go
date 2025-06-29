package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// Variables to be replaced during compilation
var (
	version = "dev"
	commit  = "unknown"
	date    = time.Now().Format(time.RFC3339)
	builtBy = "local"
)

func main() {
	// Get the Scotter project directory
	projectDir, err := getProjectDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding project directory: %s\n", err)
		os.Exit(1)
	}

	// Get the current commit if not defined
	if commit == "unknown" {
		commit, _ = getGitCommit(projectDir)
	}

	// If builtBy is not defined, use the current user
	if builtBy == "local" {
		builtBy, _ = getCurrentUser()
	}

	// Afficher les informations de build
	fmt.Println("Building Scotter...")
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Date: %s\n", date)
	fmt.Printf("Built by: %s\n", builtBy)

	// Construire les ldflags pour injecter les informations de version
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	ldflags := fmt.Sprintf("-s -w "+
		"-X github.com/caezarr-oss/scotter/pkg/version.version=%s "+
		"-X github.com/caezarr-oss/scotter/pkg/version.commit=%s "+
		"-X github.com/caezarr-oss/scotter/pkg/version.date=%s "+
		"-X github.com/caezarr-oss/scotter/pkg/version.builtBy=%s "+
		"-X github.com/caezarr-oss/scotter/pkg/version.timestamp=%s",
		version, commit, date, builtBy, timestamp)

	// Determine the executable name based on the OS
	outputName := "scotter"
	if runtime.GOOS == "windows" {
		outputName = "scotter.exe"
	}

	// Construire la commande go build
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", outputName)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the compilation
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Build failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Build completed successfully:", outputName)
}

// getProjectDir returns the Scotter project directory
func getProjectDir() (string, error) {
	// Use GOPATH or the current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// getGitCommit retourne le hash du dernier commit Git
func getGitCommit(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "unknown", err
	}
	return string(output[:len(output)-1]), nil // Remove the newline character
}

// getCurrentUser retourne le nom de l'utilisateur courant
func getCurrentUser() (string, error) {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERNAME"), nil
	}
	return os.Getenv("USER"), nil
}
