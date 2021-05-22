package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type rootConfig struct {
	SkipLibraries bool
	Seeds         []string
	Emails        []string
	GitPath       string
	OutPutPath    string
	Obfuscate     bool
	Headless      bool
}

var (
	rootCmd = &cobra.Command{
		Use:   "repo_info_extractor",
		Short: "Extract data from a Git repository",
		Long: `Use this command to extract and upload repo data your CodersRank profile.
Example usage: repo_info_extractor path --repo_path /path/to/repo`,
	}

	RootConfig rootConfig

	emailString string
	seedsString string
	Version     string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVar(&RootConfig.SkipLibraries, "skip_libraries", false, "Turns off the library detection in order to reduce the execution time.")
	rootCmd.PersistentFlags().StringVar(&emailString, "emails", "", "Predefined emails. Example: \"alim.giray@codersrank.io,alimgiray@gmail.com\"")
	rootCmd.PersistentFlags().StringVar(&seedsString, "seeds", "", "The seed is used to find similar emails. Example: \"alimgiray, alimgiray@codersrank.io\"")
	rootCmd.PersistentFlags().StringVar(&RootConfig.GitPath, "git_path", "", "where the Git binary is")
	rootCmd.PersistentFlags().StringVar(&RootConfig.OutPutPath, "output_path", "./repo_data_v2", "Where to put output file")
	rootCmd.PersistentFlags().BoolVar(&RootConfig.Obfuscate, "obfuscate", true, "File names and emails won't be hashed. Set it to true for debug purposes.")
	rootCmd.PersistentFlags().BoolVar(&RootConfig.Headless, "headless", false, "Headless mode is used on CodersRank's backend system.")
}

func initConfig() {
	emails := make([]string, 0)
	if len(emailString) > 0 {
		emails = strings.Split(emailString, ",")
	}
	RootConfig.Emails = emails

	seeds := make([]string, 0)
	if len(seedsString) > 0 {
		seeds = strings.Split(seedsString, ",")
	}

	RootConfig.Seeds = seeds
}
