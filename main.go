package main

import (
	"flag"

	"github.com/codersrank-org/repo_info_extractor/extractor"
)

func main() {

	repoPath := flag.String("repo_path", "", "Path of the repo")
	headless := flag.String("headless", "false", "Path of the repo")
	flag.Parse()

	if repoPath == nil || *repoPath == "" {
		panic("Please provide a path to the repo")
	}

	repoExtractor := extractor.RepoExtractor{
		RepoPath: *repoPath,
		Headless: *headless == "true",
	}

	err := repoExtractor.Extract()
	if err != nil {
		panic(err)
	}
}
