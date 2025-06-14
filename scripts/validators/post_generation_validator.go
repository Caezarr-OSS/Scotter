package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run post_generation_validator.go <project_dir>")
		os.Exit(1)
	}

	projectDir := os.Args[1]
	errorsFound := false
	
	// Valider les workflows GitHub générés
	workflowsDir := filepath.Join(projectDir, ".github", "workflows")
	if _, err := os.Stat(workflowsDir); err == nil {
		fmt.Printf("Validation des workflows GitHub dans: %s\n", workflowsDir)
		
		// Parcourir tous les fichiers YAML dans le répertoire workflows
		err = filepath.Walk(workflowsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// Vérifier si c'est un fichier YAML
			if !info.IsDir() && (strings.HasSuffix(strings.ToLower(path), ".yml") || strings.HasSuffix(strings.ToLower(path), ".yaml")) {
				fmt.Printf("Validation du fichier: %s\n", path)
				if err := validateYAMLFile(path); err != nil {
					fmt.Printf("❌ Fichier YAML invalide: %s - %v\n", path, err)
					errorsFound = true
					return nil // Continuer avec les autres fichiers
				}
				fmt.Printf("✅ Fichier YAML valide: %s\n", path)
			}
			
			return nil
		})
		
		if err != nil {
			fmt.Printf("Erreur lors de la validation des workflows: %v\n", err)
			os.Exit(1)
		}
	}
	
	// Valider les fichiers GoReleaser générés
	goreleaserPath := filepath.Join(projectDir, ".goreleaser.yml")
	if _, err := os.Stat(goreleaserPath); err == nil {
		fmt.Printf("Validation du fichier GoReleaser: %s\n", goreleaserPath)
		if err := validateYAMLFile(goreleaserPath); err != nil {
			fmt.Printf("❌ Fichier YAML invalide: %s - %v\n", goreleaserPath, err)
			errorsFound = true
		} else {
			fmt.Printf("✅ Fichier YAML valide: %s\n", goreleaserPath)
		}
	}

	// Vérifier si des erreurs ont été trouvées
	if errorsFound {
		fmt.Println("\n⚠️ Des erreurs de validation YAML ont été détectées. Veuillez les corriger avant de continuer.")
		os.Exit(1)
	}
	
	fmt.Println("\n✅ Tous les fichiers YAML générés sont valides.")
}

// validateYAMLFile vérifie si un fichier YAML est valide
func validateYAMLFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier: %w", err)
	}
	
	var out interface{}
	err = yaml.Unmarshal(content, &out)
	if err != nil {
		return fmt.Errorf("erreur de syntaxe YAML: %w", err)
	}
	
	return nil
}
