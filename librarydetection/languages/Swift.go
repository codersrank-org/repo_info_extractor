package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewSwiftAnalyzer constructor
func NewSwiftAnalyzer() librarydetection.Analyzer {
	return &swiftAnalyzer{}
}

type swiftAnalyzer struct{}

func (a *swiftAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	// regex to find imports like import Stuff
	regexImport, err := regexp.Compile(`import ([a-zA-Z\.]*)[\n|\r\n]`)
	if err != nil {
		return nil, err
	}
	// regex to find imports like import kind module.symbol
	// here kind can be func, var, let, typealias, protocol, enum, class, struct
	// this regex will return a list of (kind, module.symbol)
	regexDeclarations, err := regexp.Compile(`import [func|var|let|typealias|protocol|enum|class|struct]+ ([a-zA-Z|\.]*)`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{regexImport, regexDeclarations}), nil
}
