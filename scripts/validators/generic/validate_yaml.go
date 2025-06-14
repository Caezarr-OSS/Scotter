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
		fmt.Println("Usage: go run validate_yaml.go <yaml_file_or_dir>")
		os.Exit(1)
	}

	path := os.Args[1]
	isDir := false
	
	// Vérifier si le chemin est un fichier ou un répertoire
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Erreur lors de l'accès au chemin %s: %v\n", path, err)
		os.Exit(1)
	}
	
	isDir = fileInfo.IsDir()
	
	if isDir {
		fmt.Printf("Recherche de fichiers YAML dans: %s\n", path)
		// Parcourir le répertoire pour trouver tous les fichiers YAML
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// Afficher tous les fichiers trouvés
			if !info.IsDir() {
				fmt.Printf("Fichier trouvé: %s\n", path)
			}
			
			// Vérifier si c'est un fichier YAML ou un template YAML (.tmpl)
			if !info.IsDir() && (strings.HasSuffix(strings.ToLower(path), ".yml") || 
			                     strings.HasSuffix(strings.ToLower(path), ".yaml") ||
			                     (strings.HasSuffix(strings.ToLower(path), ".tmpl") && 
			                      (strings.Contains(strings.ToLower(path), ".yml.") || 
			                       strings.Contains(strings.ToLower(path), ".yaml.")))) {
				fmt.Printf("Validation du fichier: %s\n", path)
				if err := validateYAMLGeneric(path); err != nil {
					fmt.Printf("❌ Fichier YAML invalide: %s - %v\n", path, err)
					return err
				}
				fmt.Printf("✅ Fichier YAML valide: %s\n", path)
			}
			
			return nil
		})
		
		if err != nil {
			fmt.Printf("Erreur lors du traitement des fichiers: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Valider un seul fichier
		if err := validateYAMLGeneric(path); err != nil {
			fmt.Printf("❌ Fichier YAML invalide: %s - %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Fichier YAML valide: %s\n", path)
	}
}

// validateYAMLGeneric vérifie si un fichier YAML est valide
func validateYAMLGeneric(filePath string) error {
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
