# Multi-Language Support in Scotter

This document explains how Scotter supports multiple programming languages and how to extend it to support additional languages.

## Supported Languages

Scotter currently supports the following programming languages:

| Language | Type | Description |
|----------|------|-------------|
| Go | `model.GoLang` | Go projects with various types (Default, Library, CLI, API) |
| Shell | `model.NoLang` | Shell script projects without a specific language |

## Language-Specific Configuration

Each language has its own specific configuration options in the `Config` structure:

```go
type Config struct {
    ProjectName string
    Language    LanguageType
    Go          GoConfig       // Go-specific configuration
    Pipeline    PipelineConfig // Pipeline configuration
    Directories []string
}
```

### Go Configuration

Go projects have specific configuration options:

```go
type GoConfig struct {
    ProjectType GoProjectType
    ModulePath  string
}
```

The `GoProjectType` can be one of:
- `DefaultGoType`: A minimal Go project
- `LibraryGoType`: A Go library/package
- `CLIGoType`: A command-line application
- `APIGoType`: A web API/service

## Using Multiple Languages

As a user of Scotter, you can select your preferred programming language during project initialization:

```bash
scotter init
# You will be prompted to select a language (Go, Shell, etc.)
```

Based on your selection, Scotter will generate appropriate project structures and configurations.

## For Contributors: Adding a New Language

> **Note**: This section is for developers who want to contribute to the Scotter project itself by adding support for new programming languages. End users of Scotter do not need to modify any code.

To add support for a new programming language in Scotter, follow these steps:

### 1. Update the Language Type

Add a new language type constant in `internal/model/config.go`:

```go
const (
    GoLang   LanguageType = "go"
    NoLang   LanguageType = "none"
    PythonLang LanguageType = "python" // New language
)
```

### 2. Add Language-Specific Configuration

Create a new configuration structure for the language:

```go
type PythonConfig struct {
    ProjectType PythonProjectType
    PackageName string
}

type PythonProjectType string

const (
    DefaultPythonType  PythonProjectType = "default"
    LibraryPythonType  PythonProjectType = "library"
    CLIPythonType      PythonProjectType = "cli"
    WebPythonType      PythonProjectType = "web"
)
```

Update the main `Config` structure to include the new language configuration:

```go
type Config struct {
    ProjectName string
    Language    LanguageType
    Go          GoConfig
    Python      PythonConfig  // New language configuration
    Pipeline    PipelineConfig
    Directories []string
}
```

### 3. Update the Validation Logic

Update the `Validate` method to validate the new language configuration:

```go
func (cfg *Config) Validate() error {
    // Check required fields
    if cfg.ProjectName == "" {
        return fmt.Errorf("project name cannot be empty")
    }

    // Validate language type
    validLanguage := false
    for _, lang := range []LanguageType{GoLang, NoLang, PythonLang} {
        if cfg.Language == lang {
            validLanguage = true
            break
        }
    }

    if !validLanguage {
        return fmt.Errorf("invalid language type: %s", cfg.Language)
    }

    // Validate language-specific configuration
    switch cfg.Language {
    case GoLang:
        // Go validation logic
    case PythonLang:
        // Python validation logic
    }

    return nil
}
```

### 4. Create Templates

Create templates for the new language in `internal/templates/<language>/`:

```
internal/templates/python/
├── default/
│   └── main.py.tmpl
├── library/
│   └── library.py.tmpl
├── cli/
│   └── cli.py.tmpl
└── web/
    └── app.py.tmpl
```

### 5. Create a Generator

Create a generator for the new language in `internal/generator/<language>/`:

```go
package python

import (
    "github.com/Caezarr-OSS/Scotter/internal/model"
)

type Generator struct {
    Config *model.Config
}

func NewGenerator(cfg *model.Config) *Generator {
    return &Generator{
        Config: cfg,
    }
}

func (g *Generator) Generate() error {
    // Generate Python project files based on the configuration
    return nil
}
```

### 6. Update the Prompt System

Update the prompt system to allow selecting the new language:

```go
func (p *ProjectPrompt) AskLanguage() (model.LanguageType, error) {
    options := []string{
        "Go",
        "Shell Script",
        "Python", // New language
    }

    selected, err := p.AskSelect("Select the programming language for your project:", options, 0)
    if err != nil {
        return "", err
    }

    switch selected {
    case "Go":
        return model.GoLang, nil
    case "Shell Script":
        return model.NoLang, nil
    case "Python":
        return model.PythonLang, nil
    default:
        return model.GoLang, nil
    }
}
```

## Pipeline Features by Language

Different languages may require different pipeline features. Here's how pipeline features are organized by language:

| Feature | Go | Shell | Description |
|---------|-------|-------|-------------|
| `ci` | ✅ | ✅ | Continuous Integration |
| `commit-lint` | ✅ | ✅ | Conventional Commits validation |
| `changelog` | ✅ | ✅ | Changelog generation |
| `release` | ✅ | ✅ | Automated releases |
| `docker` | ✅ | ✅ | Docker container support |
| `taskfile` | ✅ | ✅ | Task automation |
| `makefile` | ✅ | ✅ | Make build system |
| `dependabot` | ✅ | ✅ | Dependency updates |

## Best Practices for Multi-Language Support

When extending Scotter to support new languages, follow these best practices:

1. **Modularity**: Keep language-specific code isolated in its own package
2. **Configuration**: Use language-specific configuration structures
3. **Templates**: Organize templates by language and project type
4. **Validation**: Add proper validation for language-specific options
5. **Documentation**: Update documentation to reflect new language support
6. **Testing**: Add tests for the new language support

## Language-Specific Directory Structures

Each language has its own recommended directory structure:

### Go Projects

```
project/
├── cmd/
│   └── project/
│       └── main.go
├── internal/
│   └── ...
├── pkg/
│   └── ...
├── go.mod
└── go.sum
```

### Shell Projects

```
project/
├── bin/
│   └── script.sh
├── lib/
│   └── functions.sh
└── README.md
```

### Python Projects (Example)

```
project/
├── project/
│   ├── __init__.py
│   └── main.py
├── tests/
│   └── test_main.py
├── setup.py
└── requirements.txt
```
