package github

import (
	"fmt"
	"os"
	"path/filepath"
)

// DebugTemplatePaths affiche des informations sur les chemins de templates
func DebugTemplatePaths() {
	// Chemins relatifs à vérifier
	paths := []string{
		"templates/github/github/ci.yml.tmpl",
		"templates/github/github/changelog.yml.tmpl",
		"templates/github/github/commitlint.yml.tmpl",
		"templates/github/github/release.yml.tmpl",
	}

	// Vérifier chaque chemin
	for _, path := range paths {
		absPath, _ := filepath.Abs(path)
		_, err := os.Stat(path)
		fmt.Printf("Path: %s\nAbsolute: %s\nExists: %v\nError: %v\n\n", path, absPath, err == nil, err)
	}

	// Vérifier également les chemins absolus
	cwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)
	
	// Vérifier les templates dans le répertoire courant
	templatesDir := filepath.Join(cwd, "templates", "github", "github")
	if _, err := os.Stat(templatesDir); err == nil {
		fmt.Printf("Templates directory exists: %s\n", templatesDir)
		files, _ := os.ReadDir(templatesDir)
		fmt.Printf("Files in templates directory:\n")
		for _, file := range files {
			fmt.Printf("- %s\n", file.Name())
		}
	} else {
		fmt.Printf("Templates directory does not exist: %s (Error: %v)\n", templatesDir, err)
	}
}
