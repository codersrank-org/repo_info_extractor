package cmd

import (
	"fmt"

	repoSource "github.com/codersrank-org/repo_info_extractor/v2/repoSources"
	"github.com/spf13/cobra"
)

type bitbucketConfig struct {
	Username   string
	Password   string
	Visibility string
}

var (
	bitbucketCmd = &cobra.Command{
		Use:   "bitbucket",
		Short: "Extract repository from BitBucket",
		Long:  `Provide the username and password and it will extract the repositories from BitBucket`,
		Run: func(cmd *cobra.Command, args []string) {
			source := repoSource.NewBitbucketProvider(
				BitbucketConfig.Username,
				BitbucketConfig.Password,
				BitbucketConfig.Visibility,
				*RootConfig.GitPath,
			)
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

	BitbucketConfig bitbucketConfig
)

func init() {
	rootCmd.AddCommand(bitbucketCmd)
	bitbucketCmd.Flags().StringVar(&BitbucketConfig.Username, "username", "", "Username to authenticate to BitBucket")
	bitbucketCmd.Flags().StringVar(&BitbucketConfig.Password, "password", "", "Password to authenticate to BitBucket")
	bitbucketCmd.Flags().StringVar(&BitbucketConfig.Visibility, "visibility", "", "Filter extracted repos by visibility. Possible values: public, private, all")
	bitbucketCmd.MarkFlagRequired("username")
	bitbucketCmd.MarkFlagRequired("password")
	bitbucketCmd.MarkFlagRequired("visibility")
}
