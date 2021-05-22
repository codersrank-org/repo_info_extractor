package repoSource

import "github.com/codersrank-org/repo_info_extractor/entities"

// RepoSource describes the interface that each provider has to implement
type RepoSource interface {
	// GetRepos provides the list of the repositories from the given provider
	GetRepos() []*entities.Repository
	// GetCLIConfigs provides the requires configuration for the given RepoSource
	// It might be a token, password and username pair, path, etc.
	GetCLIConfigs() []*entities.Config
}
