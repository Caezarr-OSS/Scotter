# Versioning Conventions in Scotter

## Tag Naming Rules

To ensure optimal compatibility with Go ecosystem tools and conventions, follow these tag naming rules:

### Go Libraries (Go Modules)

**MUST USE** a "v" prefix before the version number:

```
v0.1.0, v1.0.0, v2.0.0
```

Reason: Go dependency management tools (`go get`, `go mod`) expect library versions to follow this convention for dependency resolution.

### Non-Library Projects

Projects that are not Go libraries (CLI, API, applications):

**SHOULD NOT USE** a "v" prefix:

```
0.1.0, 1.0.0, 2.0.0
```

These projects include:
- CLI applications
- API services
- Standard applications
- Projects with local builds

## GoReleaser Configuration

Ensure your `.goreleaser.yml` file is configured to respect these conventions:

- For Go libraries, use tags with a "v" prefix
- For other projects, use tags without a "v" prefix

## GitHub Actions Workflow Files

Release generation workflows (such as `release.yml` or `go-library-release.yml`) should be configured according to the project type to use the appropriate tag format.
