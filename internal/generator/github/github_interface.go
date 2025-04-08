package github

// GitHubGenerator defines the interface for GitHub workflow generation
type GitHubGenerator interface {
	Generate() error
	GenerateCIWorkflow() error
	GenerateCommitLintWorkflow() error
	GenerateChangelogWorkflow() error
	GenerateReleaseWorkflow() error
	GenerateDependabotConfig() error
}

// Ensure Generator implements GitHubGenerator
var _ GitHubGenerator = (*Generator)(nil)
