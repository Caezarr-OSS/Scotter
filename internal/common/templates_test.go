package common

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

func TestTemplateManager(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Créer un fichier template de test
	testTmplPath := filepath.Join(tmpDir, "test.go.tmpl")
	testTmplContent := `package main

func main() {
	// {{ .ProjectName }} Project
	println("Hello, {{ .ProjectName }}!")
}
`
	if err := os.WriteFile(testTmplPath, []byte(testTmplContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Créer une configuration de test
	cfg := &model.Config{
		ProjectName: "TestProject",
		Language:    model.GoLang,
		Pipeline: model.PipelineConfig{
			CIType: model.GithubActionsCI,
		},
	}

	// Créer le TemplateManager
	tmplMgr := NewTemplateManager(cfg, tmpDir, LangGo)

	// Test 1: Vérifier la localisation des templates
	t.Run("TemplateLocator", func(t *testing.T) {
		path, err := tmplMgr.FindTemplate("test", []string{".go.tmpl"})
		if err != nil {
			t.Errorf("FindTemplate failed: %v", err)
		}
		if filepath.Base(path) != "test.go.tmpl" {
			t.Errorf("Expected test.go.tmpl, got %s", filepath.Base(path))
		}
	})

	// Test 2: Vérifier l'exécution des templates
	t.Run("TemplateExecution", func(t *testing.T) {
		// Générer un fichier à partir du template
		outputPath := filepath.Join(tmpDir, "output.go")
		err := tmplMgr.GenerateFileFromTemplate("test", outputPath, []string{".go.tmpl"})
		if err != nil {
			t.Fatalf("Template execution failed: %v", err)
		}

		// Vérifier le contenu du fichier généré
		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		expectedContent := `package main

func main() {
	// TestProject Project
	println("Hello, TestProject!")
}
`
		if string(content) != expectedContent {
			t.Errorf("Generated content doesn't match.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
		}
	})

	// Test 3: Vérifier les délimiteurs personnalisés
	t.Run("CustomDelimiters", func(t *testing.T) {
		// Définir les délimiteurs personnalisés qu'on utilisera
		customDelims := CustomDelimiters()
		
		// Vérifier que les délimiteurs sont bien ceux attendus
		if customDelims.Left != "[[" || customDelims.Right != "]]" {
			t.Fatalf("Expected custom delimiters to be [[ and ]], got %s and %s", customDelims.Left, customDelims.Right)
		}
		
		// Créer un TemplateManager avec délimiteurs personnalisés
		customTmplMgr := NewTemplateManager(cfg, tmpDir, "").WithDelimiters(customDelims)

		// Créer un contenu de template avec les délimiteurs personnalisés
		customTmplContent := fmt.Sprintf("on:\n  push:\n    branches: [\"%s .ProjectName %s\"]\n", customDelims.Left, customDelims.Right)

		// Générer un fichier
		outputPath := filepath.Join(tmpDir, "custom_output.yml")
		
		// Exécuter le template avec le contenu
		result, err := customTmplMgr.ExecuteTemplateWithContent(customTmplContent, "custom_template")
		if err != nil {
			t.Fatalf("Custom delimiter template execution failed: %v", err)
		}
		
		if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
			t.Fatalf("Failed to write output file: %v", err)
		}
		
		// Vérifier le contenu généré
		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read custom output file: %v", err)
		}

		expectedContent := "on:\n  push:\n    branches: [\"TestProject\"]\n"
		if string(content) != expectedContent {
			t.Errorf("Custom delimiter content doesn't match.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
		}
		
		fmt.Printf("Generated %s with custom delimiters\n", outputPath)
	})
}

func TestTemplateManagerIntegration(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "template-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Créer un dossier templates et des sous-dossiers
	templatesDir := filepath.Join(tmpDir, "templates")
	goTemplatesDir := filepath.Join(templatesDir, "go")
	
	if err := os.MkdirAll(goTemplatesDir, 0755); err != nil {
		t.Fatalf("Failed to create template directories: %v", err)
	}

	// Créer quelques templates
	mainTmplPath := filepath.Join(goTemplatesDir, "main.go.tmpl")
	mainTmplContent := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("{{ .ProjectName }} Project")
}
`
	if err := os.WriteFile(mainTmplPath, []byte(mainTmplContent), 0644); err != nil {
		t.Fatalf("Failed to write main template: %v", err)
	}

	// Créer une configuration
	cfg := &model.Config{
		ProjectName: "IntegrationTest",
		Language:    model.GoLang,
		Pipeline: model.PipelineConfig{
			CIType: model.GithubActionsCI,
		},
	}

	// Créer le TemplateManager
	tmplMgr := NewTemplateManager(cfg, templatesDir, LangGo)

	// Tester la génération d'un fichier
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	outputPath := filepath.Join(outputDir, "main.go")
	if err := tmplMgr.GenerateFileFromTemplate("go/main", outputPath, []string{".go.tmpl"}); err != nil {
		t.Fatalf("Failed to generate file from template: %v", err)
	}

	// Vérifier le contenu
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedContent := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("IntegrationTest Project")
}
`
	if string(content) != expectedContent {
		t.Errorf("Generated content doesn't match.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}
