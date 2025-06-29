# Scotter - Scaffolding Otter

Scotter is a powerful scaffolding tool for Go projects that allows rapid generation of project structures with integrated CI/CD workflows. It provides a modular, plugin-based architecture for easy extensibility and maintenance.

## Features

- Create Go projects with predefined templates:
  - CLI applications (using Cobra)
  - API services (using Gin)
  - Libraries
  - Default minimal structure
- Integrated CI/CD workflows:
  - GitHub Actions workflows for building, testing, and releasing
  - Conventional commits validation with commitlint
  - Multi-platform, multi-architecture support (Linux, macOS, Windows)
- GoReleaser integration for automated releases:
  - SBOM (Software Bill of Materials) generation
  - Checksums for integrity verification
  - Archive generation (zip, tar.gz)

## Installation

```bash
go install github.com/caezarr-oss/scotter@latest
```

Or build from source:

```bash
# Clone the repository
git clone https://github.com/caezarr-oss/scotter.git
cd scotter

# Install Task if you don't have it
go install github.com/go-task/task/v3/cmd/task@latest

# Build and install
task install
```

## Usage

### Initialize a new project

```bash
scotter init my-project --type cli --language go
```

Supported project types:
- `cli`: Command line application (uses Cobra)
- `api`: REST API service (uses Gin)
- `library`: Reusable library
- `default`: Minimal structure

### Add CI workflows

```bash
cd my-project
scotter add ci github
```

### Add platforms

```bash
scotter add platform linux
scotter add platform darwin
scotter add platform windows
```

### Add release assets

```bash
scotter add release-asset checksum
scotter add release-asset sbom
scotter add release-asset archive
```

## Project Configuration

Scotter uses a `.scotter.yaml` file in the project root to store configuration:

```yaml
project_name: "my-project"
project_type: "cli"
language: "go"
platforms:
  - "linux"
  - "darwin"
  - "windows"
architectures:
  - "amd64"
  - "arm64"
release_assets:
  - "checksum"
  - "sbom"
  - "archive"
ci_provider: "github"
```

## Architecture

Scotter uses a modular architecture based on interfaces to make the system extensible:

- `LanguageProvider`: Interface for language plugins
- `CIProvider`: Interface for CI/CD integrations
- `TemplateManager`: Interface for managing templates

For more details, see [ARCHITECTURE.en.md](ARCHITECTURE.en.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.