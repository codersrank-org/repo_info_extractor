package languages

import (
	"github.com/codersrank-org/repo_info_extractor/librarydetection"
	"regexp"
)

// NewGoAnalyzer constructor
func NewGoAnalyzer() librarydetection.Analyzer {
	return &goAnalyzer{}
}

type goAnalyzer struct {
}

func (a *goAnalyzer) ExtractLibraries(contents string) []string {
	// regex for multiline imports
	regex1, err := regexp.Compile(`(?msi)import\s*\(\s*(.*?)\s*\)`)
	if err != nil {
		panic(err)
	}


	// Find libraries in a multi line import
	regex2, err := regexp.Compile(`"(.*?)"`)
	if err != nil {
		panic(err)
	}

	allLibs := []string{}

	matches := regex1.FindAllStringSubmatch(contents, -1)
	for _, match := range matches {
		if len(match) > 1 {
			subgroup := match[1]
			for _, subgroupMatch := range regex2.FindAllStringSubmatch(subgroup, -1) {
				if len(subgroupMatch) > 1 {
					allLibs = append(allLibs, subgroupMatch[1:]...)
				}
			}
		}
	}


	// regex for imports like this: import _ "github.com/user/repo/..."
	regex3, err := regexp.Compile(`(?i)import[\t ]*(?:[_.].*)?[\t ]?\(?"(.+)"\)?;?\s`)
	if err != nil {
		panic(err)
	}

	//// regex for imports with alias. Like this: import alias1 "github.com/user/repo/..."
	regex4, err := regexp.Compile(`(?i)import[\t ]*[a-z].+[\t ]?\(?"(.+)"\)?;?\s`)
	if err != nil {
		panic(err)
	}

	regexes := []*regexp.Regexp{
		regex3,
		regex4,
	}

	allLibs = append(allLibs, executeRegexes(contents, regexes)...)

	return allLibs
}
