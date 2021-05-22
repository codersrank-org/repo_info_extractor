package cmd

import (
	"fmt"

	"github.com/codersrank-org/repo_info_extractor/extractor"
	"github.com/spf13/cobra"
)

type extractConfig struct {
	ReoPath  string
	RepoName string
}

var (
	pathCmd = &cobra.Command{
		Use:   "path",
		Short: "Extract repository by path",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			repoExtractor := extractor.RepoExtractor{
				RepoPath:            ExtractConfig.ReoPath,
				OutputPath:          RootConfig.OutPutPath,
				GitPath:             RootConfig.GitPath,
				Headless:            RootConfig.Headless,
				Obfuscate:           RootConfig.Obfuscate,
				UserEmails:          RootConfig.Emails,
				Seed:                RootConfig.Seeds,
				ShowProgressBar:     RootConfig.Headless != true, // Show progress bar only if running in interactive mode
				OverwrittenRepoName: ExtractConfig.RepoName,
				SkipLibraries:       RootConfig.SkipLibraries,
			}

			err := repoExtractor.Extract()
			if err != nil {
				fmt.Println("Error during execution.", err.Error())
			}
		},
	}

	ExtractConfig extractConfig
)

func init() {
	rootCmd.AddCommand(pathCmd)
	pathCmd.Flags().StringVar(&ExtractConfig.ReoPath, "repo_path", "", "Path of the repo")
	pathCmd.Flags().StringVar(&ExtractConfig.RepoName, "repo_name", "", "You can overwrite the default repo name. This name will be shown on the profile page.")
}
