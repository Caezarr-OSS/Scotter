# Scotter Templates Documentation

> **Note**: This document is primarily for developers who want to contribute to the Scotter project itself. End users of Scotter do not need to modify templates or understand the internal template system.

This document explains how to use, customize, and extend the templates used by Scotter to generate projects in multiple programming languages.

## Template Overview

Scotter uses Go's built-in `text/template` package to generate project files from templates. These templates are located in the `internal/templates` directory and are organized by language and feature:

```
internal/templates/
├── common/                  # Language-agnostic templates
│   ├── readme.md.tmpl       # README template
│   └── taskfile.yml.tmpl    # Taskfile template
├── github/                  # GitHub-related templates
│   ├── ci.yml.tmpl          # GitHub Actions CI workflow
│   ├── commitlint.yml.tmpl  # Commitlint workflow 
│   └── release.yml.tmpl     # Release workflow using GoReleaser
├── go/                      # Go-specific templates
│   ├── default_main.go.tmpl # Template for minimal project main.go
│   ├── library/             # Library-specific templates
│   ├── cli/                 # CLI-specific templates
│   └── api/                 # API-specific templates
└── shell/                   # Shell-specific templates
    ├── script.sh.tmpl       # Unix shell script template
    └── script.ps1.tmpl      # PowerShell script template
```

## Template Variables

Templates can access project configuration values using the following variables:

### Common Variables

| Variable | Description | Example |
|----------|-------------|----------|
| `.ProjectName` | Name of the project | `myproject` |
| `.Language` | Programming language | `model.GoLang` |
| `.Directories` | Directories to create | `["cmd", "internal", "pkg"]` |

### Language-Specific Variables

#### Go Variables

| Variable | Description | Example |
|----------|-------------|----------|
| `.Go.ProjectType` | Type of Go project | `model.CLIGoType` |
| `.Go.ModulePath` | Go module path | `github.com/username/myproject` |

#### Python Variables

| Variable | Description | Example |
|----------|-------------|----------|
| `.Python.ProjectType` | Type of Python project | `model.LibraryPythonType` |
| `.Python.PackageName` | Python package name | `myproject` |

### Pipeline Configuration

| Variable | Description | Example |
|----------|-------------|----------|
| `.Pipeline.UseGitHubActions` | Whether to use GitHub Actions | `true` |
| `.Pipeline.SelectedFeatures` | List of enabled pipeline features | `["ci", "commit-lint", "changelog"]` |

## Customizing Existing Templates

You can customize how Scotter generates files by modifying the templates:

1. Locate the template you want to modify in the `internal/templates` directory
2. Edit the template using Go's template syntax
3. Test your changes by running Scotter

### Example: Customizing the README template

The `readme.md.tmpl` template is used to generate the project's README. Here's how you might customize it:

```go
# {{.ProjectName}}

{{.Description}}

## Features

- Feature 1
- Feature 2
{{if contains .Pipeline.SelectedFeatures "taskfile"}}
- Includes task automation with Taskfile
{{end}}

{{if eq .Language "model.GoLang"}}
## Installation

```bash
go install {{.Go.ModulePath}}@latest
```

## Usage

```bash
{{.ProjectName}} --help
```
{{end}}
```

This example shows how to use conditional statements based on language and selected features. The `contains` function checks if a feature is in the list of selected features.

## Adding New Templates

To add a new template:

1. Create a new `.tmpl` file in the appropriate directory (language-specific or common)
2. Update the corresponding generator to use your new template
3. Modify the configuration model if needed to support new options

### Example: Adding Support for a New Language

1. Update the `model/config.go` file to add a new language type:

```go
const (
	GoLang LanguageType = "go"
	NoLang LanguageType = "none"
	PythonLang LanguageType = "python"  // New language
)
```

2. Create a directory for the new language's templates:

```
mkdir -p internal/templates/python
```

3. Add language-specific templates, for example `internal/templates/python/main.py.tmpl`:

```python
#!/usr/bin/env python3

def main():
    print("Hello from {{.ProjectName}}!")

if __name__ == "__main__":
    main()
```

4. Create a new generator in `internal/generator/python/python.go`

5. Update the prompt system to allow selecting the new language

### Example: Adding a New Pipeline Feature

1. Update the `model/config.go` file to add the new feature to the pipeline features:

```go
var AvailablePipelineFeatures = []PipelineFeature{
	// Existing features
	{Name: "ci", Description: "Continuous Integration", Dependencies: nil},
	{Name: "commit-lint", Description: "Conventional Commits validation", Dependencies: nil},
	// New feature
	{Name: "docker", Description: "Docker container support", Dependencies: nil},
}
```

2. Create templates for the new feature, for example `internal/templates/docker/Dockerfile.tmpl`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
{{if eq .Language "model.GoLang"}}
RUN go mod download
RUN go build -o /bin/app ./cmd/{{.ProjectName}}
{{end}}

FROM alpine:latest
{{if eq .Language "model.GoLang"}}
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
{{end}}
```

3. Create a new generator in `internal/generator/docker/docker.go`
3. Update the initializer to use the new generator

## Cross-Platform Considerations

When creating or modifying templates, consider cross-platform compatibility:

- Use `{{if eq .OS "windows"}}` for Windows-specific sections
- Use path manipulation that works on all platforms
- Be careful with line endings in shell scripts

## Template Function Reference

Beyond the standard Go template functions, Scotter provides the following custom functions:

| Function | Description | Example |
|----------|-------------|---------|
| `toLower` | Convert string to lowercase | `{{toLower .ProjectName}}` |
| `toUpper` | Convert string to uppercase | `{{toUpper .ProjectName}}` |
| `toTitle` | Convert string to title case | `{{toTitle .ProjectName}}` |

## Testing Templates

Test your templates by running Scotter with different configurations:

```bash
go run ./cmd/scotter/main.go
```

Check the generated files to ensure they match your expectations.
