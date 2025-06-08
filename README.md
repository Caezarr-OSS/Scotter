# Scotter

Scotter is a flexible scaffolding tool that simplifies project setup, GitHub environment configuration, and development best practices integration. Originally focused on Go, Scotter now supports multiple programming languages and offers modular pipeline features.

![Scotter Logo](assets/img/scotter.png)

## Features

### Multi-Language Support
- **Go**: Supports various project types (Default, Library, CLI, API)
- **Shell**: For script-based projects
- More languages coming soon!

### Modular Pipeline Features
- **GitHub Actions**: Customizable CI/CD workflows
- **Commit Validation**: Enforces conventional commits with commitlint
- **Automatic Releases**: Generates releases with changelogs
- **Build Tools**: Sets up Taskfile and/or Makefile based on project needs
- **Container Support**: Generates Dockerfile or Containerfile based on your preference

### Project Structure
- Organizes files according to language-specific best practices
- Creates sensible defaults while allowing customization
- Generates comprehensive documentation

## Installation

```bash
go install github.com/Caezarr-OSS/Scotter/cmd/scotter@latest
```

## Usage

### Project Initialization

In a new Git repository, run with the required command line arguments:

```bash
scotter init --name myproject --go-type cli --features ci,commit-lint,changelog,release \
  --os linux,darwin --arch amd64,arm64
```

At minimum, you need to specify the project name with `--name`.

Key flags include:
- `--name`: Project name
- `--lang`: Programming language (go, none)
- `--go-type`: Go project type (default, library, cli, api)
- `--features`: Pipeline features (comma-separated)
- `--os`: Target operating systems (comma-separated)
- `--arch`: Target architectures (comma-separated)

Use `scotter init --help` to see all available options.

### Build Target Management

Scotter allows adding, removing, and listing build targets for your Go project at any time:

```bash
# List current build targets
scotter target list

# Add a new build target
scotter target add --os darwin --arch arm64

# Remove a build target
scotter target remove --os windows --arch amd64
```

This enables flexible cross-compilation configuration as your project evolves.

## Project Types

### Go Projects
- **Default/Minimal**: Simple structure for scripting or generic projects
- **Library**: For reusable Go packages
- **CLI**: For command-line tools
- **API/Service**: For web servers/APIs

### Shell Projects
- **Basic**: Simple shell script structure with best practices

### Pipeline Features
- **CI**: Continuous integration with tests and linting
- **Commit Lint**: Enforces conventional commit messages
- **Changelog**: Automatic changelog generation
- **Release**: Automated release process
- **Documentation**: Generates and maintains documentation
- **Container**: Creates container configuration files (Dockerfile or Containerfile)

### GitHub Release Configuration

To enable GitHub releases through workflows, you'll need to configure a personal access token with appropriate permissions:

1. **Create a GitHub personal access token**:
   - For classic tokens: Select the full `repo` scope
   - For fine-grained tokens: Add repository permissions for `Contents: Read and write`, `Metadata: Read-only`, and `Actions: Read and write`

2. **Add the token to your repository secrets**:
   - Go to your repository → Settings → Secrets and variables → Actions
   - Create a new repository secret named `RELEASE_TOKEN`
   - Paste your personal access token as the value

3. **The release workflow will automatically use this token** for creating releases, uploading assets, and publishing changelogs

## Development

Scotter follows the GitFlow model with a `develop` branch for development and a `main` branch for releases.

### Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using conventional commits:
   ```bash
   git commit -m 'feat(scope): add some amazing feature'
   ```
   Valid types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`
   
   Valid scopes: `core`, `model`, `prompt`, `generator`, `config`, `init`, `cli`, `docs`, `deps`
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Tools

Scotter uses the same tools it provides to its users:

- **Task Runner**: Use `task` to run development tasks
  ```bash
  # Run tests
  task test
  
  # Run linters
  task lint
  
  # Build the project
  task build
  
  # Check commit messages
  task commitlint
  
  # Update changelog
  task changelog
  ```

- **Conventional Commits**: All commits must follow the [Conventional Commits](https://www.conventionalcommits.org/) specification

- **Automated Changelog**: The CHANGELOG.md is automatically generated from commit messages

## Templates

Scotter utilise un système de templates standardisé pour générer le contenu des projets. Tous les templates sont situés dans le répertoire `/templates` à la racine du projet.

### Organisation des Templates

- Templates racine : Templates généraux utilisés pour tous les types de projets
- `/templates/github` : Templates spécifiques pour les workflows GitHub Actions
- `/templates/container` : Templates pour la génération de fichiers Docker/Containerfile

### Types de Templates

- **Fichiers de projet** : Templates pour les fichiers sources (main.go, example.go, etc.)
- **Documentation** : Templates pour README.md, API.md, etc.
- **CI/CD** : Templates pour les workflows d'intégration et de déploiement continus
- **Builds** : Templates pour les outils de build comme Taskfile, Makefile, etc.

### Customizing Templates

All templates use Go's text/template syntax and can be customized as needed. The template system is designed to be extensible, allowing users to modify existing templates or create new ones for specific project requirements.

## Documentation

### User Documentation
- [User Guide](docs/USER_GUIDE.md)
- [FAQ](docs/FAQ.md)
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md)

### Feature Documentation
- [Multi-Language Support](docs/MULTI_LANGUAGE_SUPPORT.md)
- [Container Support](docs/CONTAINER_SUPPORT.md)
- [Container Best Practices](docs/CONTAINER_BEST_PRACTICES.md)
- [Changelog Support](docs/CHANGELOG_SUPPORT.md)

### Developer Documentation
- [Developer Guide](docs/DEVELOPER_GUIDE.md)
- [Contributing Guide](docs/CONTRIBUTING.md)

## License

Distributed under the MIT License. See the `LICENSE` file for more information.
