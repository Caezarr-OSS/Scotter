# Scotter Templates Documentation

This document explains how to use, customize, and extend the templates used by Scotter to generate Go projects.

## Template Overview

Scotter uses Go's built-in `text/template` package to generate project files from templates. These templates are located in the `internal/templates` directory and are organized as follows:

```
internal/templates/
├── default_main.go.tmpl     # Template for minimal project main.go
├── github/                  # GitHub-related templates
│   ├── ci.yml.tmpl          # GitHub Actions CI workflow
│   ├── commitlint.yml.tmpl  # Commitlint workflow 
│   └── release.yml.tmpl     # Release workflow using GoReleaser
├── readme.md.tmpl           # README template
└── taskfile.yml.tmpl        # Taskfile template
```

## Template Variables

Templates can access project configuration values using the following variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `.ProjectName` | Name of the project | `myproject` |
| `.ProjectType` | Type of the project | `model.CLIType` |
| `.ModulePath` | Go module path | `github.com/username/myproject` |
| `.Features.UseTaskFile` | Whether to include Taskfile | `true` |
| `.Features.GitHub.UseWorkflows` | Whether to include GitHub workflows | `true` |
| `.Features.GitHub.UseCommitLint` | Whether to use commitlint | `true` |
| `.Features.GitHub.UseReleaseWorkflow` | Whether to use release workflow | `true` |
| `.Features.GitHub.UseDependabot` | Whether to use Dependabot | `true` |
| `.Features.GitHub.GenerateChangelog` | Whether to generate CHANGELOG | `true` |

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
{{if .Features.UseTaskFile}}
- Includes task automation with Taskfile
{{end}}

## Installation

```bash
go install {{.ModulePath}}@latest
```

## Usage

```bash
{{.ProjectName}} --help
```
```

This example shows how to use conditional statements (`{{if .Features.UseTaskFile}}`) to include content based on project configuration.

## Adding New Templates

To add a new template:

1. Create a new `.tmpl` file in the appropriate directory
2. Update the corresponding generator to use your new template
3. Modify the configuration model if needed to support new options

### Example: Adding a Docker template

1. Create a new file `internal/templates/docker/Dockerfile.tmpl`:

```
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /bin/app ./cmd/{{.ProjectName}}

FROM alpine:latest
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
```

2. Create a new generator in `internal/generator/docker/docker.go`
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
