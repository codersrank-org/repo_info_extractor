package repoSource

import "github.com/codersrank-org/repo_info_extractor/entities"

type directoryPath struct{}

func NewDirectoryPath() RepoSource {
	return &directoryPath{}
}

func (dp *directoryPath) GetRepos() []*entities.Repository {
	return nil
}

func (dp *directoryPath) GetCLIConfigs() []*entities.Config {
	return nil
}
