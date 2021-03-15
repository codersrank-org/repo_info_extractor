package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewCAnalyzer constructor
func NewCAnalyzer() librarydetection.Analyzer {
	return &cAnalyzer{}
}

type cAnalyzer struct {}

func (a *cAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	regex, err := regexp.Compile(`(?i)#include\s?[<"]([/a-zA-Z0-9.-]+)[">]`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{regex}), nil
}
