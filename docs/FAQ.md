# Frequently Asked Questions (FAQ)

## General Questions

### What is Scotter?

Scotter is a flexible scaffolding tool that simplifies project setup, GitHub environment configuration, and development best practices integration. It supports multiple programming languages and offers modular pipeline features.

### Why should I use Scotter instead of other scaffolding tools?

Scotter differentiates itself by:
- Supporting multiple programming languages with language-specific best practices
- Offering modular pipeline features that can be enabled/disabled based on your needs
- Providing container support with both Dockerfile and Containerfile formats
- Integrating seamlessly with GitHub Actions and other CI/CD tools
- Following a consistent, opinionated approach to project structure

### Which programming languages does Scotter support?

Currently, Scotter supports:
- Go (with multiple project types: Default, Library, CLI, API)
- Shell (for script-based projects)

More languages are planned for future releases.

## Installation and Setup

### How do I install Scotter?

```bash
go install github.com/Caezarr-OSS/Scotter/cmd/scotter@latest
```

### What are the prerequisites for using Scotter?

- Go 1.18 or later
- Git
- A terminal/command prompt

### Can I use Scotter in an existing project?

Yes, but use caution. Scotter is primarily designed for initializing new projects. When using it with existing projects:
1. Make sure to back up your project first
2. Use the `--force` flag if you want to overwrite existing files
3. Consider initializing in a temporary directory and then manually merging the generated files

## Features and Usage

### What pipeline features does Scotter offer?

Scotter offers several modular pipeline features:
- GitHub Actions for CI/CD
- Commit validation with commitlint
- Automatic changelog generation
- Release automation
- Documentation generation
- Container support (Dockerfile/Containerfile)

### How do I select which features to enable?

During the initialization process, Scotter will prompt you to select which features you want to enable. You can use the arrow keys to navigate and the space bar to select/deselect features.

### Can I add features to an existing Scotter project?

Yes, you can run `scotter init` again with the `--force` flag to add new features. However, this might overwrite customizations you've made to existing files.

## Container Support

### What's the difference between Dockerfile and Containerfile?

The difference is primarily in the naming convention:
- `Dockerfile` is the standard name used by Docker
- `Containerfile` is the standard name used by Podman and other OCI-compliant tools

The content and syntax of both files are identical and follow the same OCI specification.

### Which container format should I choose?

Choose based on the container engine you primarily use:
- Select `Dockerfile` if you primarily use Docker
- Select `Containerfile` if you primarily use Podman or other OCI-compliant tools

### Can I switch between Dockerfile and Containerfile later?

Yes, you can simply rename the file. Since both formats use identical syntax, renaming from `Dockerfile` to `Containerfile` (or vice versa) is all you need to do.

### Does Scotter optimize container images for my language?

Yes, Scotter generates optimized container configurations based on your selected programming language:
- Go projects use multi-stage builds to create smaller images
- Shell projects use Alpine Linux with appropriate script configurations
- Other project types get a sensible default configuration

## GitHub Actions Integration

### What GitHub Actions workflows does Scotter generate?

Depending on your selected features, Scotter can generate workflows for:
- CI (testing and linting)
- Release automation
- Container image building and publishing

### Do I need to set up any secrets for the GitHub Actions workflows?

Yes, for some features:
- For releases, you'll need a `GITHUB_TOKEN` (automatically provided by GitHub)
- For container publishing, you'll need appropriate registry credentials

### Can I customize the GitHub Actions workflows?

Yes, the generated workflow files are just starting points. You can edit them to suit your specific needs.

## Troubleshooting

### Scotter isn't generating all the files I expected

Make sure you've selected all the features you want during initialization. You can run `scotter init` again with the `--force` flag to regenerate files.

### The generated container file doesn't work with my project

The generated container files are templates based on common patterns. You may need to customize them for your specific project requirements. See our [Container Best Practices](CONTAINER_BEST_PRACTICES.md) guide for tips.

### GitHub Actions workflows are failing

Check the workflow logs on GitHub for specific errors. Common issues include:
- Missing secrets
- Incorrect file paths
- Dependency issues

See our [Troubleshooting Guide](TROUBLESHOOTING.md) for more help.

## Contributing and Support

### How can I contribute to Scotter?

See our [Contributing Guide](CONTRIBUTING.md) for details on how to contribute to Scotter.

### Where can I report bugs or request features?

Please create an issue on our [GitHub repository](https://github.com/Caezarr-OSS/Scotter/issues).

### Is there a community or forum for Scotter users?

You can use the [Discussions](https://github.com/Caezarr-OSS/Scotter/discussions) section of our GitHub repository to connect with other users and ask questions.
