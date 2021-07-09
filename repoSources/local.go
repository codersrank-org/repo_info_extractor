package repoSource

import (
	"github.com/codersrank-org/repo_info_extractor/v2/entities"
	"os"
	"strings"
)

type directoryPath struct {
	// path is the directory path of the repository
	path string
	// name is an optional name that can be overwritten by the user
	name string
}

// NewDirectoryPath create a new RepoSource instance.
// path is the directory path to the repository.
func NewDirectoryPath(path, name string) RepoSource {
	return &directoryPath{
		path: path,
		name: name,
	}
}

// GetRepos does nothing in this case because we already work with a local copy
func (dp *directoryPath) GetRepos() []*entities.Repository {
	fullName := ""
	if dp.name == "" {
		names := strings.Split(dp.path, string(os.PathSeparator))
		fullName = names[len(names)-1]
	} else {
		fullName = dp.name
	}

	repo := &entities.Repository{
		ID:       "",
		FullName: fullName,
		Name:     fullName,
	}
	return []*entities.Repository{repo}
}

// Clone does nothing in this case because we already work with a local copy
func (dp *directoryPath) Clone(repository *entities.Repository) (string, error) {
	return dp.path, nil
}

// CleanUp does not have to clean up anything.
func (dp *directoryPath) CleanUp() {}
