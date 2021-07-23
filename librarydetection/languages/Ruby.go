package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection"
)

// NewRubyScriptAnalyzer constructor
func NewRubyScriptAnalyzer() librarydetection.Analyzer {
	return &RubyScriptAnalyzer{}
}

type RubyScriptAnalyzer struct{}

func (a *RubyScriptAnalyzer) ExtractLibraries(contents string) ([]string, error) {

	requireRegex, err := regexp.Compile(`(?:require|gem)[( ]{1}['"]([a-zA-Z0-9_\-/]+)["'][)]?`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{requireRegex}), nil
}
