package common

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// TemplateDelimiters définit les délimiteurs personnalisés pour différents types de templates
type TemplateDelimiters struct {
	Left  string
	Right string
}

// DefaultDelimiters renvoie les délimiteurs par défaut pour les templates Go
func DefaultDelimiters() TemplateDelimiters {
	return TemplateDelimiters{
		Left:  "{{",
		Right: "}}",
	}
}

// CustomDelimiters renvoie les délimiteurs personnalisés (utilisés pour GitHub Actions, YAML, etc.)
func CustomDelimiters() TemplateDelimiters {
	return TemplateDelimiters{
		Left:  "{{",
		Right: "}}",
	}
}

// TemplateLanguage représente un langage supporté par le générateur
type TemplateLanguage string

const (
	// Langages supportés
	LangGo       TemplateLanguage = "go"
	LangPython   TemplateLanguage = "python"
	LangRust     TemplateLanguage = "rust"
	LangBash     TemplateLanguage = "bash"
	LangTypeScript TemplateLanguage = "typescript"
)

// TemplateManager gère l'exécution des templates pour tous les générateurs
type TemplateManager struct {
	Config      *model.Config
	BaseTmplDir string
	Language    TemplateLanguage
	Delimiters  TemplateDelimiters
}

// NewTemplateManager crée un nouveau gestionnaire de templates
func NewTemplateManager(config *model.Config, baseTmplDir string, language TemplateLanguage) *TemplateManager {
	return &TemplateManager{
		Config:      config,
		BaseTmplDir: baseTmplDir,
		Language:    language,
		Delimiters:  DefaultDelimiters(),
	}
}

// WithDelimiters définit les délimiteurs personnalisés pour ce gestionnaire
func (m *TemplateManager) WithDelimiters(delimiters TemplateDelimiters) *TemplateManager {
	m.Delimiters = delimiters
	return m
}

// FindTemplate cherche un template dans plusieurs chemins possibles
func (m *TemplateManager) FindTemplate(templateName string, extensions []string) (string, error) {
	fmt.Printf("Finding template: %s with extensions %v in base dir %s\n", templateName, extensions, m.BaseTmplDir)
	// Les chemins possibles pour le template, par ordre de priorité:
	// 1. templates/{templateName}_{language}{.ext} - Spécifique au langage
	// 2. templates/{templateName}{.ext} - Générique

	possiblePaths := []string{}
	
	// Si un langage est spécifié, chercher d'abord les templates spécifiques au langage
	if m.Language != "" {
		for _, ext := range extensions {
			// Format: templates/{templateName}_{language}{.ext}
			langPath := filepath.Join(m.BaseTmplDir, templateName+"_"+string(m.Language)+ext)
			possiblePaths = append(possiblePaths, langPath)
		}
	}
	
	// Ensuite ajouter les chemins génériques
	for _, ext := range extensions {
		genericPath := filepath.Join(m.BaseTmplDir, templateName+ext)
		possiblePaths = append(possiblePaths, genericPath)
	}
	
	// Chercher les templates dans l'ordre
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	
	// Aucun template trouvé
	return "", fmt.Errorf("no template found for %s (looked in: %v)", templateName, possiblePaths)
}

// ExecuteTemplate exécute un template et retourne le résultat
func (m *TemplateManager) ExecuteTemplate(templatePath string) (string, error) {
	// Lire le template
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}
	
	return m.ExecuteTemplateWithContent(string(tmplContent), templatePath)
}

// ExecuteTemplateWithContent exécute un template à partir de son contenu
func (m *TemplateManager) ExecuteTemplateWithContent(content string, name string) (string, error) {
	// Créer un template avec les délimiteurs appropriés
	tmpl := template.New(filepath.Base(name)).Delims(m.Delimiters.Left, m.Delimiters.Right)
	
	// Compiler le template
	parsedTmpl, err := tmpl.Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}
	
	// Exécuter le template
	var buf bytes.Buffer
	if err := parsedTmpl.Execute(&buf, m.Config); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}
	
	return buf.String(), nil
}

// GenerateFileFromTemplate génère un fichier à partir d'un template
func (m *TemplateManager) GenerateFileFromTemplate(templateName, outputPath string, extensions []string) error {
	// Chercher le template
	templatePath, err := m.FindTemplate(templateName, extensions)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}
	
	// Exécuter le template
	output, err := m.ExecuteTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	
	// Créer le répertoire de destination si nécessaire
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
	}
	
	// Écrire le résultat dans le fichier
	if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}
	
	fmt.Printf("Generated %s from %s\n", outputPath, filepath.Base(templatePath))
	return nil
}

// GenerateFileFromContent génère un fichier à partir d'un contenu prédéfini (sans template)
func (m *TemplateManager) GenerateFileFromContent(outputPath string, content string) error {
	// Créer le répertoire de destination si nécessaire
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", outputPath, err)
	}
	
	// Écrire le contenu dans le fichier
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}
	
	fmt.Printf("Generated %s\n", outputPath)
	return nil
}
