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

func (a *cAnalyzer) ExtractLibraries(contents string) []string {
	regex, err := regexp.Compile(`#include\s?[<"]([/a-zA-Z0-9.-]+)[">]`)
	if err != nil {
		panic(err)
	}
	matches := regex.FindAllStringSubmatch(contents, -1)

	var allLibs []string
	for _, match := range matches {
		if len(match) > 1 {
			allLibs = append(allLibs, match[1:]...)
		}
	}


	return allLibs
}
