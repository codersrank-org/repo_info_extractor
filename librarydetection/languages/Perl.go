package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewPerlAnalyzer constructor
func NewPerlAnalyzer() librarydetection.Analyzer {
	return &perlAnalyzer{}
}

type perlAnalyzer struct{}

func (a *perlAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	regex, err := regexp.Compile(`(?:use|require)[^\S\n]+(?:if.*,\s+)?[\"']?([a-zA-Z][a-zA-Z0-9:]*)[\"']?(?:\s+.*)?;`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{regex}), nil
}
