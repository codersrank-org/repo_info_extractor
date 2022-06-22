package languages

import (
	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection"
	"regexp"
)

// NewRustAnalyzer constructor
func NewRustAnalyzer() librarydetection.Analyzer {
	return &rustAnalyzer{}
}

type rustAnalyzer struct {
}

func (a *rustAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	// regex to find "uses", braces are not supported as their tree syntax is too complicated for Go's regexes
	useRegex, err := regexp.Compile(`use\s+(?:::)?(?:r#)?([^;:\s]+)`)
	if err != nil {
		return nil, err
	}

	// regex for "extern crate" statements
	externCrateRegex, err := regexp.Compile(`extern\s+crate\s+(?:r#)?(\S+)\s*;`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{useRegex, externCrateRegex}), nil
}
