package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewCAnalyzer constructor
func NewCppAnalyzer() librarydetection.Analyzer {
	return &cppAnalyzer{}
}

type cppAnalyzer struct {}

func (a *cppAnalyzer) ExtractLibraries(contents string) []string {
	regex, err := regexp.Compile(`(?i)#include\s?[<"]([/a-zA-Z0-9.-]+)[">]`)
	if err != nil {
		panic(err)
	}

	return executeRegexes(contents, []*regexp.Regexp{regex})
}
