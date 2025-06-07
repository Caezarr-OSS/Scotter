package model

// CIType définit le type de système CI/CD
type CIType string

const (
	// NoneCI pour désactiver la génération de CI
	NoneCI CIType = "none"
	
	// GithubActionsCI pour les workflows GitHub Actions (nom utilisé dans le code)
	GithubActionsCI CIType = "github"
	// GitHubActions pour les workflows GitHub Actions (nom alternatif)
	GitHubActions CIType = "github"
	
	// GitlabCI pour les pipelines GitLab CI (nom utilisé dans le code)
	GitlabCI CIType = "gitlab"
	// GitLabCI pour les pipelines GitLab CI (nom alternatif)
	GitLabCI CIType = "gitlab"
	
	// CircleCI pour les configurations CircleCI
	CircleCI CIType = "circleci"
	
	// TravisCI pour les configurations Travis CI
	TravisCI CIType = "travis"
)

// String renvoie la représentation en string du CIType
func (c CIType) String() string {
	return string(c)
}

// AllCITypes renvoie tous les types de CI supportés
func AllCITypes() []CIType {
	return []CIType{
		GithubActionsCI,
		GitlabCI,
		CircleCI,
		TravisCI,
	}
}
