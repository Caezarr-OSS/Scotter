// +build ignore

// This script copies template files from internal/templates to internal/embedded/templates
// to be embedded in the binary using go:embed
package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	sourceDir := "internal/templates"
	targetDir := "internal/embedded/templates"

	// Clear target directory first
	if err := os.RemoveAll(targetDir); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error removing target directory: %v\n", err)
		os.Exit(1)
	}

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating target directory: %v\n", err)
		os.Exit(1)
	}

	// Walk through the source directory and copy all files
	err := filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Skip source directory itself
		if rel == "." {
			return nil
		}

		// Target path
		targetPath := filepath.Join(targetDir, rel)

		// Handle directories
		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Handle files
		return copyFile(path, targetPath)
	})

	if err != nil {
		fmt.Printf("Error copying templates: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Templates successfully copied to embedded directory")
}

func copyFile(src, dst string) error {
	// Open source file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy contents
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Sync to ensure file is written
	return out.Sync()
}
