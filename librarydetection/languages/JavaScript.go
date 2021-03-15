package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewJavaScriptAnalyzer constructor
func NewJavaScriptAnalyzer() librarydetection.Analyzer {
	return &javaScriptAnalyzer{}
}

type javaScriptAnalyzer struct{}

func (a *javaScriptAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	require, err := regexp.Compile(`require\(["\'](.+)["\']\);?\s`)
	importRegex, err := regexp.Compile(`import\s*(?:.+ from)?\s?\(?[\'"](.+)[\'"]\)?;?\s`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{require, importRegex}), nil
}
