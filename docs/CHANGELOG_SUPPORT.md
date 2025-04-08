# Changelog Support in Scotter

This document explains how Scotter's automatic changelog generation feature works and how to use it effectively in your projects.

## Overview

The changelog feature in Scotter automatically generates and maintains a structured changelog based on your project's commit history. It follows the [Keep a Changelog](https://keepachangelog.com/) format and [Semantic Versioning](https://semver.org/) principles.

## How It Works

When you select the changelog feature during project initialization, Scotter:

1. Adds a GitHub Actions workflow (`.github/workflows/changelog.yml`) that automatically generates and updates the CHANGELOG.md file
2. Adds a task in the Taskfile.yml for local changelog generation
3. Configures the necessary dependencies and tools

The workflow:
- Triggers on pushes to the main/master branch
- Can be manually triggered via GitHub Actions interface
- Ignores changes to the CHANGELOG.md file itself to prevent infinite loops
- Automatically commits and pushes the updated changelog

## Requirements

To use the changelog feature effectively:

- Your project must use [Conventional Commits](https://www.conventionalcommits.org/) format for commit messages (enforced by the Commit Lint feature)
- Node.js must be available in your CI environment (automatically installed in GitHub Actions)
- For local generation, Node.js and npx must be installed on your machine

## GitHub Actions Workflow

The generated workflow file (`.github/workflows/changelog.yml`) looks like this:

```yaml
name: Generate Changelog

on:
  # Manual trigger
  workflow_dispatch:
  # Automatic trigger on push to main/master with conventional commits
  push:
    branches:
      - main
      - master
    paths-ignore:
      - 'CHANGELOG.md'

jobs:
  generate-changelog:
    name: Generate Changelog
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 16
      
      - name: Install conventional-changelog-cli
        run: npm install -g conventional-changelog-cli
      
      - name: Generate changelog
        run: conventional-changelog -p angular -i CHANGELOG.md -s
      
      - name: Commit and push if changed
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add CHANGELOG.md
          git diff --quiet && git diff --staged --quiet || git commit -m "docs: update changelog [skip ci]"
          git push
```

## Taskfile Integration

The changelog task in Taskfile.yml allows you to generate the changelog locally:

```yaml
changelog:
  desc: Generate or update CHANGELOG.md
  cmds:
    - |
      if ! command -v npx &> /dev/null; then
        echo "Error: npx not found. Please install Node.js"
        exit 1
      fi
    - |
      if ! command -v conventional-changelog &> /dev/null; then
        npm install -g conventional-changelog-cli
      fi
    - conventional-changelog -p angular -i CHANGELOG.md -s
```

To use it, simply run:

```bash
task changelog
```

## Commit Message Format

For the changelog generation to work properly, your commit messages should follow the Conventional Commits format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Common types include:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files

## Best Practices

1. **Use Conventional Commits**: Always use the conventional commit format for your commit messages
2. **Be Descriptive**: Write clear and descriptive commit messages that explain what changed and why
3. **Breaking Changes**: Mark breaking changes with `BREAKING CHANGE:` in the commit footer
4. **Scope**: Use scopes to indicate which part of the codebase was modified
5. **Generate Before Release**: Always generate the changelog before creating a new release

## Troubleshooting

If you encounter issues with changelog generation:

1. **Commit Format**: Ensure your commits follow the conventional format
2. **Node.js**: Verify that Node.js is installed and accessible
3. **Dependencies**: Make sure conventional-changelog-cli is installed
4. **Git History**: The changelog generator needs access to the full git history, ensure your checkout has sufficient depth

## Further Reading

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
- [conventional-changelog-cli](https://github.com/conventional-changelog/conventional-changelog/tree/master/packages/conventional-changelog-cli)
