package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewCAnalyzer constructor
func NewCSharpAnalyzer() librarydetection.Analyzer {
	return &cSparpAnalyzer{}
}

type cSparpAnalyzer struct {}

func (a *cSparpAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	regex1, err := regexp.Compile(`(?i)using\s?([/a-zA-z0-9.]+);`)
	if err != nil {
		return nil, err
	}

	regex2, err := regexp.Compile(`(?i)using [/a-zA-z0-9.]+ = ([/a-zA-z0-9.]+);`)
	if err != nil {
		return nil, err
	}

	regexes := []*regexp.Regexp{
		regex1,
		regex2,
	}

	return executeRegexes(contents, regexes), nil
}
