package generator

import (
	"fmt"
	"path/filepath"

	"github.com/Caezarr-OSS/Scotter/internal/common"
	"github.com/Caezarr-OSS/Scotter/internal/generator/ci"
	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// ProjectGenerator orchestre la génération d'un projet complet
type ProjectGenerator struct {
	Config      *model.Config
	TemplateDir string
	ciManager   ci.CIManager
}

// NewProjectGenerator crée un nouveau générateur de projet
func NewProjectGenerator(config *model.Config, templateDir string) (*ProjectGenerator, error) {
	// Créer la factory pour les générateurs CI
	ciFactory := ci.NewCIManagerFactory(templateDir)
	
	// Obtenir le manager CI approprié pour la configuration
	ciManager, err := ciFactory.CreateManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create CI manager: %w", err)
	}
	
	return &ProjectGenerator{
		Config:      config,
		TemplateDir: templateDir,
		ciManager:   ciManager,
	}, nil
}

// Generate génère la structure complète du projet
func (g *ProjectGenerator) Generate() error {
	fmt.Printf("Generating %s project: %s\n", g.Config.Language, g.Config.ProjectName)
	
	// Générer la structure de base du projet selon le langage
	if err := g.generateBaseStructure(); err != nil {
		return fmt.Errorf("failed to generate base structure: %w", err)
	}
	
	// Générer les fichiers spécifiques au langage
	if err := g.generateLanguageSpecificFiles(); err != nil {
		return fmt.Errorf("failed to generate language specific files: %w", err)
	}
	
	// Générer la configuration CI/CD si applicable
	if g.ciManager != nil {
		if err := g.ciManager.Generate(); err != nil {
			return fmt.Errorf("failed to generate CI/CD configuration: %w", err)
		}
	}
	
	fmt.Println("Project generation completed successfully!")
	return nil
}

// generateBaseStructure crée la structure de base du projet
func (g *ProjectGenerator) generateBaseStructure() error {
	// Créer les répertoires de base du projet
	for _, dirPath := range g.Config.Directories {
		// TODO: Implémenter la création des répertoires
		_ = dirPath // Pour éviter l'erreur de variable non utilisée
	}
	
	// Générer le README.md
	templateMgr := common.NewTemplateManager(g.Config, g.TemplateDir, "")
	if err := templateMgr.GenerateFileFromTemplate("readme", "README.md", []string{".md.tmpl"}); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}
	
	// Générer .gitignore
	if err := templateMgr.GenerateFileFromTemplate("gitignore", ".gitignore", []string{".tmpl"}); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}
	
	return nil
}

// generateLanguageSpecificFiles génère les fichiers spécifiques au langage
func (g *ProjectGenerator) generateLanguageSpecificFiles() error {
	switch g.Config.Language {
	case model.GoLang:
		return g.generateGoFiles()
	// Ajouter d'autres langages à l'avenir
	default:
		return fmt.Errorf("unsupported language: %s", g.Config.Language)
	}
}

// generateGoFiles génère les fichiers spécifiques à Go
func (g *ProjectGenerator) generateGoFiles() error {
	// Créer un TemplateManager spécifique à Go
	templateMgr := common.NewTemplateManager(g.Config, g.TemplateDir, common.TemplateLanguage("go"))
	
	// Générer go.mod
	if err := templateMgr.GenerateFileFromTemplate("go.mod", "go.mod", []string{".tmpl"}); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}
	
	// Générer le fichier main.go ou autres selon le type de projet
	switch g.Config.Go.ProjectType {
	case model.LibraryGoType:
		// Générer la structure de bibliothèque
		if err := templateMgr.GenerateFileFromTemplate("lib_main", filepath.Join(g.Config.ProjectName+".go"), []string{".tmpl", ".go.tmpl"}); err != nil {
			return fmt.Errorf("failed to generate library main file: %w", err)
		}
	case model.CLIGoType:
		// Générer la structure CLI
		if err := templateMgr.GenerateFileFromTemplate("cli_main", filepath.Join("cmd", g.Config.ProjectName, "main.go"), []string{".tmpl", ".go.tmpl"}); err != nil {
			return fmt.Errorf("failed to generate CLI main file: %w", err)
		}
	// Autres types de projets...
	}
	
	return nil
}
