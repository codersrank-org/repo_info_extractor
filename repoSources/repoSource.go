package repoSource

import (
	"fmt"
	"io/ioutil"

	"github.com/codersrank-org/repo_info_extractor/entities"
	"github.com/codersrank-org/repo_info_extractor/extractor"
)

type ExtractConfig struct {
	OutputPath      string
	GitPath         string
	Headless        bool
	Obfuscate       bool
	UserEmails      []string
	Seeds           []string
	ShowProgressBar bool // Show progress bar only if running in interactive mode
	SkipLibraries   bool
}

// RepoSource describes the interface that each provider has to implement
type RepoSource interface {
	// GetRepos provides the list of the repositories from the given provider
	GetRepos() []*entities.Repository
	// Clone clones the given repository to the given directory
	// returns with the cloned path and an error if any.
	Clone(repository *entities.Repository) (string, error)
	// CleanUp revert to the state before extraction.
	// E.g. remove temporary files, directories.
	CleanUp()
}

func ExtractFromSource(source RepoSource, config ExtractConfig) error {
	repos := source.GetRepos()

	if config.OutputPath == "" {
		outputDir, err := ioutil.TempDir("", "clone_dir_")
		if err != nil {
			return fmt.Errorf("couldn't create temp dir for artifacts. Try to set it with --output_path. Error: %s", err.Error())
		}
		config.OutputPath = outputDir
	}

	for _, repo := range repos {
		path, err := source.Clone(repo)
		if err != nil {
			fmt.Println("Couldn't clone repository. Error:", err.Error())
		}

		repoExtractor := extractor.RepoExtractor{
			RepoPath:            path,
			OutputPath:          config.OutputPath + "/" + repo.GetSafeFullName(),
			GitPath:             config.GitPath,
			Headless:            config.Headless,
			Obfuscate:           config.Obfuscate,
			UserEmails:          config.UserEmails,
			Seed:                config.Seeds,
			ShowProgressBar:     config.Headless != true, // Show progress bar only if running in interactive mode
			OverwrittenRepoName: repo.Name,
			SkipLibraries:       config.SkipLibraries,
		}

		err = repoExtractor.Extract()
		if err != nil {
			fmt.Println("Error during execution.", err.Error())
			continue
		}

	}

	artifactUploader := NewArtifactUploader(config.OutputPath)
	artifactUploader.UploadRepos(repos)

	source.CleanUp()

	return nil
}
