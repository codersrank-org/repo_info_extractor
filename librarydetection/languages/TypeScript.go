package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection"
)

// NewTypeScriptAnalyzer constructor
func NewTypeScriptAnalyzer() librarydetection.Analyzer {
	return &typeScriptAnalyzer{}
}

type typeScriptAnalyzer struct{}

func (a *typeScriptAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	require, err := regexp.Compile(`require\(["\'](.+)["\']\);?\s`)
	if err != nil {
		return nil, err
	}
	importRegex, err := regexp.Compile(`import\s*(?:.+ from)?\s?\(?[\'"](.+)[\'"]\)?;?\s`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{require, importRegex}), nil
}
