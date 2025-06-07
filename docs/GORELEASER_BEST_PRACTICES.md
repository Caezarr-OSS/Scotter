# GoReleaser Best Practices

This document outlines best practices for using GoReleaser with Scotter projects.

## Configuration File Structure

A properly configured `.goreleaser.yml` should include these key sections:

```yaml
# GoReleaser configuration
project_name: myproject
before:
  hooks:
    - go mod tidy
    - go generate ./...
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.Version={{.Version}}
      - -X main.CommitSHA={{.Commit}}
```

## Template Variables

GoReleaser provides variables that can be used in the configuration:

| Variable | Description | Example |
|----------|-------------|---------|
| `.Version` | The version being released | `v1.2.3` |
| `.Tag` | The tag name | `v1.2.3` |
| `.Commit` | The commit SHA | `abcdef123456` |
| `.Date` | The release date | `2025-01-02_15:04:05` |
| `.Env` | Access environment variables | `{{.Env.HOME}}` |

## Common Pitfalls and Solutions

### 1. Template Variable Syntax

✅ **Correct**:
```yaml
ldflags:
  - -X main.Version={{.Version}}
```

❌ **Incorrect**:
```yaml
ldflags:
  - -X main.Version=${VERSION}
```

### 2. Archive Names

✅ **Correct**:
```yaml
archives:
  - id: default
    name_template: >-
      {{ .ProjectName }}_
      {{- .Version }}_
      {{- .Os }}_
      {{- .Arch }}
```

❌ **Incorrect**:
```yaml
archives:
  - id: default
    name_template: {{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}
```

### 3. Changelog Generation

For comprehensive changelog generation, always use the `git` option:

```yaml
changelog:
  use: git
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - Merge pull request
      - Merge branch
```

## Required Files

GoReleaser may require these files in your repository:

1. **LICENSE**: Required for proper license attribution in packages
2. **README.md**: Used for documentation in package managers

## GitHub Workflow Integration

When using GoReleaser with GitHub Actions, configure proper permissions:

```yaml
jobs:
  goreleaser:
    permissions:
      contents: write
    steps:
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v4
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.RELEASE_TOKEN }}
```

## Testing GoReleaser Locally

Before pushing, test your configuration:

```bash
# Validate config
goreleaser check

# Test release (without publishing)
goreleaser release --snapshot --clean --skip=publish
```

## SBOM Generation

For Software Bill of Materials generation:

```yaml
sboms:
  - artifacts: archive
```

## References

- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [GoReleaser Config Templates](https://goreleaser.com/customization/templates/)
