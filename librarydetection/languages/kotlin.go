package languages

import (
	"regexp"
	"strings"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection"
)

// NewCAnalyzer constructor
func NewKotlinAnalyzer() librarydetection.Analyzer {
	return &kotlinAnalyzer{}
}

type kotlinAnalyzer struct{}

func (a *kotlinAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	// regex to find imports like org.example (exclude standard java kotlin libraries)
	regex1, err := regexp.Compile(`(?i)import ([a-zA-Z0-9.]*[^.*\n])`)
	if err != nil {
		return nil, err
	}

	ret := executeRegexes(contents, []*regexp.Regexp{regex1})
	var res = []string{}
	for _, v := range ret {
		// remove all those starting with `kotlin` or `java`
		if strings.HasPrefix(v, "java") || strings.HasPrefix(v, "kotlin") {
			continue
		}

		res = append(res, v)
	}

	return res, nil
}
