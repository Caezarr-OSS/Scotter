# Troubleshooting Guide for Scotter

This guide helps you solve common issues you might encounter when using Scotter.

## Installation Issues

### "Command not found" after installation

If you see `scotter: command not found` after installation:

1. Make sure your Go bin directory is in your PATH:
   ```bash
   echo $PATH | grep -q "$(go env GOPATH)/bin" || echo "$(go env GOPATH)/bin is not in your PATH"
   ```

2. Add it to your PATH if needed:
   ```bash
   # For Bash/Zsh
   echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
   source ~/.bashrc
   
   # For Fish
   fish_add_path (go env GOPATH)/bin
   ```

3. Verify installation:
   ```bash
   which scotter
   ```

### Permission Issues

If you encounter permission issues during installation:

```bash
sudo go install github.com/Caezarr-OSS/Scotter/cmd/scotter@latest
```

## Project Initialization Issues

### Initialization Fails

If `scotter init` fails:

1. Make sure you're in a valid Git repository:
   ```bash
   git status
   ```

2. If not, initialize a Git repository:
   ```bash
   git init
   ```

3. Try again with verbose logging:
   ```bash
   scotter init --verbose
   ```

### Template Rendering Errors

If you see template errors:

1. Make sure you have the latest version:
   ```bash
   go install github.com/Caezarr-OSS/Scotter/cmd/scotter@latest
   ```

2. Check if your project name contains special characters that might cause template issues

## Container Configuration Issues

### Container Build Fails

If your container build fails:

1. Verify the Dockerfile/Containerfile was generated correctly:
   ```bash
   cat Dockerfile # or Containerfile
   ```

2. Make sure Docker/Podman is installed and running:
   ```bash
   docker --version # or podman --version
   ```

3. Try building manually to see detailed errors:
   ```bash
   docker build -t myproject . # or podman build -t myproject .
   ```

### GitHub Actions Container Workflow Fails

If the container GitHub Action workflow fails:

1. Make sure your GitHub repository has the necessary secrets set up for container registry access
2. Check that the container file path in the workflow matches your project structure
3. Verify that the container registry you're using is accessible from GitHub Actions

## GitHub Actions Issues

### Workflow Files Not Generated

If GitHub Actions workflow files aren't generated:

1. Make sure you selected GitHub Actions during initialization
2. Check if the `.github/workflows` directory exists
3. Run initialization again with the `--force` flag to regenerate files

### Workflows Fail on GitHub

If workflows fail when running on GitHub:

1. Check the workflow logs for specific errors
2. Verify that your project builds and tests pass locally
3. Make sure any required secrets are configured in your repository settings

## Language-Specific Issues

### Go Projects

If you encounter issues with Go projects:

1. Verify your Go version meets the minimum requirements:
   ```bash
   go version
   ```

2. Make sure your go.mod and go.sum files are correctly generated:
   ```bash
   cat go.mod
   ```

3. Try running `go mod tidy` to fix dependency issues

### Shell Projects

If you encounter issues with Shell projects:

1. Check file permissions on generated scripts:
   ```bash
   ls -la bin/
   ```

2. Make scripts executable if needed:
   ```bash
   chmod +x bin/*.sh
   ```

## Getting More Help

If you're still experiencing issues:

1. Check the [GitHub Issues](https://github.com/Caezarr-OSS/Scotter/issues) to see if someone has reported the same problem
2. Open a new issue with detailed information about your problem
3. Include logs, error messages, and steps to reproduce the issue
