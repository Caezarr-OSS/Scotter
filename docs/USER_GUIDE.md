# Scotter User Guide

This guide explains how to use Scotter to create projects in different programming languages and configure CI/CD pipelines for them.

## Getting Started

### Installation

Install Scotter using Go:

```bash
go install github.com/Caezarr-OSS/Scotter/cmd/scotter@latest
```

### Basic Usage

To initialize a new project:

```bash
scotter init
```

This will start an interactive prompt that guides you through the project setup.

## Selecting a Programming Language

When you run `scotter init`, you will be prompted to select a programming language for your project:

```
? Select the programming language for your project:
  ▸ Go
    Shell Script
```

Currently, Scotter supports:
- **Go**: For Go projects (libraries, CLI tools, APIs)
- **Shell Script**: For shell script projects

Select your preferred language using the arrow keys and press Enter.

## Configuring Project Type

After selecting a language, you'll be prompted to choose a project type. The available options depend on the language you selected:

### Go Project Types

```
? Select the type of Go project:
  ▸ Default/Minimal (Simple structure)
    Library (Reusable package)
    CLI (Command-line tool)
    API/Service (Web server)
```

Each project type creates a different directory structure and includes different files:

- **Default/Minimal**: Basic Go project structure
- **Library**: Structure optimized for reusable Go packages
- **CLI**: Structure for command-line applications with flags and commands
- **API/Service**: Structure for web servers/APIs with routing and middleware

### Shell Project Types

```
? Select the type of shell project:
  ▸ Basic (Simple script)
```

## Configuring Pipeline Features

After selecting the project type, you'll be prompted to configure pipeline features:

```
? Do you want to use GitHub Actions for CI/CD? (Y/n)
```

If you choose to use GitHub Actions, you'll be prompted to select which pipeline features to enable:

```
? Select pipeline features to enable:
  ▸ CI (Continuous Integration)
    Commit Lint (Conventional Commits validation)
    Changelog (Automatic changelog generation)
    Release (Automated releases)
    Docker (Container support)
```

You can select multiple features using the space bar, and then press Enter to confirm.

### Available Pipeline Features

| Feature | Description | What it does |
|---------|-------------|--------------|
| CI | Continuous Integration | Sets up testing and linting on each commit |
| Commit Lint | Conventional Commits validation | Enforces commit message format |
| Changelog | Automatic changelog generation | Creates and updates CHANGELOG.md based on commits |
| Release | Automated releases | Creates GitHub releases and binaries |
| Containers | Container support | Adds container configuration and build workflow |

## Project Configuration

After selecting the language, project type, and pipeline features, you'll be prompted for additional configuration:

### For Go Projects

```
? Enter the project name: myproject
? Enter the Go module path: github.com/username/myproject
```

### For Shell Projects

```
? Enter the project name: myscript
```

## Generated Project Structure

Scotter will generate a project structure based on your selections. Here are examples of what you'll get:

### Go Project (CLI Type)

```
myproject/
├── .github/
│   └── workflows/
│       ├── ci.yml          # Includes build, test and commit lint jobs
│       ├── changelog.yml
│       └── release.yml     # Includes SBOM generation
├── cmd/
│   └── myproject/
│       └── main.go
├── internal/
│   └── ...
├── pkg/
│   └── ...
├── .commitlintrc.json     # Commit message convention config
├── CHANGELOG.md
├── go.mod
├── LICENSE
├── README.md
└── Taskfile.yml
```

### Shell Project (Basic Type)

```
myscript/
├── .github/
│   └── workflows/
│       ├── ci.yml          # Includes build, test and commit lint jobs
│       └── changelog.yml
├── bin/
│   └── myscript.sh
├── .commitlintrc.json     # Commit message convention config
├── CHANGELOG.md
├── LICENSE
└── README.md
```

## Using Pipeline Features

Once your project is generated, you can use the pipeline features as follows:

### Continuous Integration (CI)

The CI workflow runs automatically on each push to your repository. It performs:
- Linting
- Testing
- Building

You don't need to do anything special to trigger it.

### Commit Lint

When you enable the commit lint feature, Scotter creates a dedicated job in the CI workflow that runs independently of the build and test process. This ensures that commit message validation provides fast feedback without blocking the main build pipeline.

To make commits that pass the commit lint validation:

```bash
# Format: <type>(<scope>): <subject>
git commit -m "feat(cli): add new command for user management"
```

Valid types include: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`

The validation runs automatically on every push and pull request, ensuring all commits follow the conventional commits specification.

### Changelog Generation

To update the changelog based on your commits:

```bash
task changelog
```

This will parse your conventional commits and update CHANGELOG.md.

### Changelog Generation

When you select the Changelog feature during project initialization, Scotter sets up automatic changelog generation based on conventional commit messages.

### How It Works

1. Scotter adds a GitHub Actions workflow (`.github/workflows/changelog.yml`) that automatically generates and updates the CHANGELOG.md file
2. The workflow triggers on pushes to the main/master branch and can also be manually triggered
3. It uses `conventional-changelog-cli` to parse commit messages and generate a structured changelog
4. Changes are automatically committed and pushed back to the repository

### Requirements

- Your project must use conventional commit messages (enforced by the Commit Lint feature)
- Node.js must be available in your CI environment (automatically installed in GitHub Actions)

### Using the Changelog Feature

Once set up, the changelog will be automatically updated whenever you push commits to the main branch. The generated CHANGELOG.md follows the [Keep a Changelog](https://keepachangelog.com/) format and groups changes by type:

- Features (feat)
- Bug fixes (fix)
- Documentation changes (docs)
- And more...

You can also manually generate the changelog locally using the task provided in the Taskfile:

```bash
task changelog
```

This requires Node.js and will install the necessary dependencies if they're not already present.

## Container Support

When you select the container feature, you'll be asked which container file format you prefer:

```
? Select your preferred container file format:
  ▸ Dockerfile (Docker standard)
    Containerfile (Podman/OCI standard)
```

#### Dockerfile vs Containerfile

This choice determines the name of the generated container configuration file:

- **Dockerfile**: Standard format used by Docker
- **Containerfile**: Standard format used by Podman and other OCI-compliant container engines

It's important to understand that **the content of both files is identical** - only the filename changes to match your preferred container engine's convention. The syntax and directives used in both files follow the same specification.

#### Why Choose One Over the Other?

- Choose **Dockerfile** if you primarily use Docker or if your CI/CD pipeline expects this filename
- Choose **Containerfile** if you primarily use Podman, Buildah, or other OCI-compliant tools

#### Template Customization

Scotter generates different container configurations based on your selected programming language:

**For Go projects**:
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
# Copy go.mod and go.sum files
COPY go.mod go.sum* ./
# Download dependencies
RUN go mod download
# Copy source code
COPY . .
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/{{.ProjectName}} ./cmd/{{.ProjectName}}

# Final stage
FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/bin/{{.ProjectName}} /app/{{.ProjectName}}
# Run the application
ENTRYPOINT ["/app/{{.ProjectName}}"]
```

**For Shell projects**:
```dockerfile
FROM alpine:latest
WORKDIR /app
# Copy scripts
COPY bin/ /app/bin/
COPY . /app/
# Make scripts executable
RUN chmod +x /app/bin/*.sh
# Set the entry point to the main script
ENTRYPOINT ["/app/bin/{{.ProjectName}}.sh"]
```

**For other project types**:
```dockerfile
FROM alpine:latest
WORKDIR /app
# Copy all files
COPY . /app/
# Set the entry point to a shell
CMD ["/bin/sh"]
```

Scotter will also generate a GitHub Actions workflow for building and publishing container images, which works with both Dockerfile and Containerfile formats.

### Releases

To create a new release:

1. Tag your commit:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The release workflow will automatically:
   - Build binaries for multiple platforms
   - Generate Software Bill of Materials (SBOM) files for each release artifact
   - Create a GitHub release
   - Attach the binaries, SBOMs, and checksums to the release
   - Update the changelog

### Software Bill of Materials (SBOM)

All projects generated by Scotter include automatic SBOM generation as part of the release process. This provides transparency about your software's components and dependencies.

- **For applications**: GoReleaser generates SPDX-format SBOMs using Syft for each release artifact
- **For libraries**: A dedicated step in the release workflow generates an SBOM for the package

These SBOM files can be used by security scanning tools to identify vulnerabilities and ensure compliance with security policies.

## Customizing Your Project

After generating your project, you can customize it further:

1. Edit the README.md to add more information about your project
2. Modify the generated code to implement your functionality
3. Add more files and directories as needed

## Using Scotter for Existing Projects

To add Scotter features to an existing project:

1. Navigate to your project directory
2. Run `scotter init`
3. Follow the prompts to configure the project
4. Scotter will add the necessary files without overwriting your existing code

## Troubleshooting

### Common Issues

1. **Git not initialized**: Scotter requires a Git repository. If you get an error, run:
   ```bash
   git init
   ```

2. **Missing dependencies**: Some features require external tools. Install them as needed:
   ```bash
   # For commit linting (used by GitHub Actions, not needed locally with .commitlintrc.json)
   npm install -g @commitlint/cli @commitlint/config-conventional
   
   # For task running
   go install github.com/go-task/task/v3/cmd/task@latest
   ```

3. **GitHub Actions not running**: Make sure your repository is connected to GitHub and you have enabled Actions in the repository settings.

## Getting Help

If you encounter any issues or have questions:

1. Check the [GitHub Issues](https://github.com/Caezarr-OSS/Scotter/issues)
2. Open a new issue if your problem isn't already reported
3. Join the community discussions
