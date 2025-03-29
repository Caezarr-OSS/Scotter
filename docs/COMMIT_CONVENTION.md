# Commit Convention Guide

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification to ensure consistent and readable commit messages.

## Commit Message Structure

Each commit message consists of a **header**, a **body**, and a **footer**:

```
<type>(<scope>): <subject>

<body>

<footer>
```

- **Header**: Required, cannot be longer than 100 characters
- **Body**: Optional, providing detailed description
- **Footer**: Optional, for referencing issues and breaking changes

## Types

The `type` field must be one of the following:

| Type | Description |
|------|-------------|
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation changes |
| `style` | Changes that do not affect the meaning of the code (formatting, etc.) |
| `refactor` | Code changes that neither fix a bug nor add a feature |
| `perf` | Code changes that improve performance |
| `test` | Adding or correcting tests |
| `build` | Changes to the build system or external dependencies |
| `ci` | Changes to CI configuration files and scripts |
| `chore` | Other changes that don't modify src or test files |
| `revert` | Reverts a previous commit |

## Scope

The `scope` field is optional and should be a noun describing the section of the codebase affected by the change:

Examples:
- `feat(generator)`
- `fix(prompt)`
- `docs(README)`
- `refactor(model)`

## Subject

The `subject` should be a short description of the change:

- Use imperative, present tense: "change" not "changed" nor "changes"
- Don't capitalize the first letter
- No period (.) at the end

## Body

The `body` should include the motivation for the change and contrast this with previous behavior:

```
feat(cli): add option to specify output format

Add a new option to the CLI that allows users to specify the format of the
output (json, yaml, or text). This makes it easier to integrate with other
tools by providing machine-readable output formats.
```

## Footer

The `footer` should contain information about Breaking Changes and reference GitHub issues:

```
BREAKING CHANGE: The API endpoint /users has been renamed to /accounts

Closes #123, #456
```

## Breaking Changes

Breaking changes should be indicated by `BREAKING CHANGE:` at the beginning of the footer or body section:

```
feat(api): change authentication mechanism

BREAKING CHANGE: The authentication token format has changed. 
All clients need to be updated to use the new format.
```

You can also use an exclamation mark after the type/scope to highlight that the change contains breaking changes:

```
feat!: change API response format
```

## Examples

### Simple Feature
```
feat(generator): add support for Docker templates
```

### Bug Fix with Issue Reference
```
fix(structure): correct path handling on Windows

Closes #42
```

### Documentation Update
```
docs: improve installation instructions
```

### Breaking Change
```
feat!(cli): change command-line arguments

BREAKING CHANGE: The --config flag has been renamed to --configuration.
```

## Tools

This project uses:

1. **commitlint**: Validates commit messages against the conventional commits format
2. **GitHub Actions**: Enforces commit message format in PRs
3. **CHANGELOG generation**: Automatically generates a changelog based on commit messages

## Using with Git Hooks

The project provides a pre-commit hook to validate your commit messages:

1. The hook is automatically installed when you run Scotter
2. It will verify your commit messages match the conventional format
3. If the message doesn't match, the commit will be rejected with guidance

## For Contributors

Please follow these commit conventions when contributing to this project. It helps with:

1. **Automatically generated changelogs**
2. **Semantic versioning** determination
3. **Navigable git history** with meaningful commit messages
4. **Easier code review** process with context-rich commits
