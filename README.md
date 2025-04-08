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

In a new Git repository, run:

```bash
scotter init
```

The tool will guide you through a series of questions to configure your project:

1. Select your programming language
2. Choose a project type based on the selected language
3. Configure pipeline features (GitHub Actions, commit validation, etc.)
4. Customize additional settings specific to your project

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
