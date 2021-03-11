package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewGoAnalyzer constructor
func NewGoAnalyzer() librarydetection.Analyzer {
	return &goAnalyzer{
		Regexes: createRegexes(),
	}
}

type goAnalyzer struct {
	Regexes []*regexp.Regexp
}

func createRegexes() []*regexp.Regexp {
	regexes := make([]*regexp.Regexp, 0, 5)

	// Find things between double quotes
	doubleQuoteRegex, err := regexp.Compile("\"(.*?)\"")
	if err == nil {
		regexes = append(regexes, doubleQuoteRegex)
	}

	return regexes
}

func (a *goAnalyzer) ExtractLibraries(contents string) []string {
	allLibs := make([]string, 0, 5)
	for _, r := range a.Regexes {
		matches := r.FindAllStringSubmatch(contents, -1)
		for _, match := range matches {
			if len(match) > 1 {
				libs := match[1:]
				allLibs = append(allLibs, libs...)
			}
		}
	}
	return allLibs
}
