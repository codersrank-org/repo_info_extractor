package main

import (
	"github.com/codersrank-org/repo_info_extractor/autoupdater"
	"github.com/codersrank-org/repo_info_extractor/cmd"
)

var (
	version = "v9.9.9" // Version of the file. E.g v0.9.6. This is set during build time.
)

func main() {
	cmd.Version = version
	au := autoupdater.NewAutoUpdater(version)
	au.CheckUpdates()

	cmd.Execute()
}
