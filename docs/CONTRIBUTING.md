# Contributing to Scotter

Thank you for your interest in contributing to Scotter! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our code of conduct. Please be respectful and considerate of others.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue in the GitHub repository with the following information:

1. A clear, descriptive title
2. A detailed description of the issue
3. Steps to reproduce the bug
4. Expected behavior
5. Actual behavior
6. Screenshots or logs (if applicable)
7. Environment information (OS, Go version, etc.)

### Suggesting Enhancements

We welcome suggestions for enhancements! Please create an issue with:

1. A clear, descriptive title
2. A detailed description of the proposed enhancement
3. Any relevant examples or mockups
4. Why this enhancement would be useful to most Scotter users

### Pull Requests

1. Fork the repository
2. Create a new branch from `develop` (not `main`)
3. Make your changes
4. Add or update tests as necessary
5. Update documentation to reflect your changes
6. Ensure all tests pass
7. Submit a pull request to the `develop` branch

#### Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that do not affect the meaning of the code (formatting, etc.)
- `refactor`: Code changes that neither fix a bug nor add a feature
- `perf`: Performance improvements
- `test`: Adding or correcting tests
- `build`: Changes to the build system or dependencies
- `ci`: Changes to CI configuration
- `chore`: Other changes that don't modify source or test files
- `revert`: Reverts a previous commit

**Scopes:**
- `core`: Core functionality
- `model`: Data models
- `prompt`: User prompts
- `generator`: File generators
- `config`: Configuration handling
- `init`: Project initialization
- `cli`: Command-line interface
- `docs`: Documentation
- `deps`: Dependencies
- `container`: Container support

**Examples:**
```
feat(container): add support for Containerfile format
fix(prompt): correct validation for project name input
docs(user-guide): update container configuration section
```

## Development Environment Setup

1. Install Go (version 1.18 or later)
2. Clone the repository
3. Install Task (taskfile.dev)
4. Run `task setup` to install development dependencies

## Testing

- Run `task test` to run all tests
- Run `task test:unit` for unit tests only
- Run `task test:integration` for integration tests only

## Code Style

We use `gofmt` and `golint` to maintain code style. Run `task lint` to check your code.

## Documentation

When adding or modifying features, please update the relevant documentation:

- `USER_GUIDE.md` for user-facing documentation
- Code comments for developer documentation
- Specific feature documentation files as appropriate

## Review Process

1. All pull requests require at least one review from a maintainer
2. CI checks must pass
3. Documentation must be updated as necessary
4. Tests must be added or updated as necessary

## Getting Help

If you need help with contributing, please:

1. Check the existing documentation
2. Look for similar issues in the issue tracker
3. Ask questions in the discussion section

Thank you for contributing to Scotter!
