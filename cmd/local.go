package cmd

import (
	"fmt"

	repoSource "github.com/codersrank-org/repo_info_extractor/repoSources"
	"github.com/spf13/cobra"
)

type extractConfig struct {
	RepoPath string
	RepoName string
}

var (
	localCmd = &cobra.Command{
		Use:   "local",
		Short: "Extract local repository by path",
		Run: func(cmd *cobra.Command, args []string) {
			source := repoSource.NewDirectoryPath(ExtractConfig.RepoPath, ExtractConfig.RepoName)
			config := repoSource.ExtractConfig{
				OutputPath:      *RootConfig.OutPutPath,
				GitPath:         *RootConfig.GitPath,
				Headless:        *RootConfig.Headless,
				Obfuscate:       *RootConfig.Obfuscate,
				UserEmails:      *RootConfig.Emails,
				Seeds:           *RootConfig.Seeds,
				ShowProgressBar: !*RootConfig.Headless,
				SkipLibraries:   *RootConfig.SkipLibraries,
			}
			err := repoSource.ExtractFromSource(source, config)

			if err != nil {
				fmt.Println("Couldn't locally extract repo. Error:", err.Error())
			}
		},
	}

	ExtractConfig extractConfig
)

func init() {
	rootCmd.AddCommand(localCmd)
	localCmd.Flags().StringVar(&ExtractConfig.RepoPath, "repo_path", "", "Path of the repo")
	localCmd.MarkFlagRequired("repo_path")
	localCmd.Flags().StringVar(&ExtractConfig.RepoName, "repo_name", "", "You can overwrite the default repo name. This name will be shown on the profile page.")
}
