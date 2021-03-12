package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewTypeScriptAnalyzer constructor
func NewTypeScriptAnalyzer() librarydetection.Analyzer {
	return &typeScriptAnalyzer{}
}

type typeScriptAnalyzer struct{}

func (a *typeScriptAnalyzer) ExtractLibraries(contents string) []string {
	require, err := regexp.Compile(`require\(["\'](.+)["\']\);?\s`)
	importRegex, err := regexp.Compile(`import\s*(?:.+ from)?\s?\(?[\'"](.+)[\'"]\)?;?\s`)
	if err != nil {
		panic(err)
	}

	return executeRegexes(contents, []*regexp.Regexp{require, importRegex})
}
