# Scotter



Scotter is a scaffolding tool for Go projects that simplifies GitHub environment setup and development best practices integration.

![Scotter Logo](assets/img/scotter.png)

## Features

- Initializes different types of Go projects (Default, Library, CLI, API, Complete)
- Configures GitHub Actions (CI/CD, commit validation, automatic releases)
- Sets up build tools (Taskfile, Makefile)
- Configures commit validation (commitlint)
- Structures the project according to Go best practices

## Installation

```bash
go install github.com/Caezarr-OSS/Scotter/cmd/scotter@latest
```

## Usage

In a new Git repository, run:

```bash
scotter init
```

The tool will ask you a series of questions to configure your project.

## Project Types

- **Default/Minimal**: Simple structure for scripting or generic projects
- **Library**: For reusable Go packages
- **CLI**: For command-line tools
- **API/Service**: For web servers/APIs
- **Complete**: All features enabled

## Development

Scotter follows the GitFlow model with a `develop` branch for development and a `main` branch for releases.

### Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See the `LICENSE` file for more information.
