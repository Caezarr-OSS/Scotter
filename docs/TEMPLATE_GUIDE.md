# Scotter Template Guide

This guide explains how templates work in Scotter and provides best practices for creating and maintaining templates.

## Template Structure

Scotter uses Go's text/template package for all templating needs. Templates are organized by feature type in the `internal/templates` directory.

### Directory Structure

```
internal/templates/
├── github/            # GitHub Actions workflow templates
├── go/                # Go project templates
│   ├── cli/           # CLI-specific templates
│   ├── api/           # API-specific templates
│   └── library/       # Library-specific templates
└── taskfile/          # Taskfile templates
```

## Custom Delimiters

Scotter uses custom delimiters `[[ ]]` instead of the default `{{ }}` for GitHub Action workflow templates to prevent conflicts with GitHub Actions' own expression syntax.

Example:
```yaml
name: CI

on:
  push:
    branches: [develop, main, 'feature/*', 'release/*']
  pull_request:
    branches: [develop, main]

jobs:
  build:
    name: Build and Test
    runs-on: ubuntu-latest
    steps:
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
```

## GoReleaser Integration

When generating GoReleaser configuration files, ensure the proper template variable syntax is used:

✅ Correct:
```yaml
ldflags:
  - -s -w
  - -X main.Version={{.Version}}
  - -X main.CommitSHA={{.Commit}}
```

❌ Incorrect:
```yaml
ldflags:
  - -s -w
  - -X main.Version=${VERSION}
  - -X main.CommitSHA=${COMMIT}
```

## Best Practices

1. **Validate Templates**: Always run tests to ensure templates can be parsed correctly
2. **Use Consistent Naming**: Follow Go conventions for template variable naming
3. **Handle Default Values**: Provide sensible defaults where possible
4. **Document Requirements**: Note any required files (e.g., LICENSE) in the template comments

## Common Issues

### GitHub Actions Template Problems

GitHub Actions uses `${{ }}` syntax for expressions. To avoid conflicts with Go templates:

1. Use `[[ ]]` delimiters for Go templates in GitHub Action workflow files
2. Ensure proper escaping when both syntaxes are needed

### GoReleaser Configuration

GoReleaser uses `{{.Version}}` style variables. Common issues include:

1. Using `${VERSION}` style variables instead of the correct syntax
2. Missing required files referenced in the configuration (e.g., LICENSE)
3. Not maintaining consistency in naming templates

## Testing Templates

Run the GitHub workflow template tests to verify functionality:

```bash
go test -v ./internal/generator/github/...
```
