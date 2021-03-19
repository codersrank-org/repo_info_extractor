package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/librarydetection"
)

// NewJavaAnalyzer constructor
func NewJavaAnalyzer() librarydetection.Analyzer {
	return &javaAnalyzer{}
}

type javaAnalyzer struct{}

func (a *javaAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	// regex to find imports like org.springframework (excluding standart java libraries)
	regex1, err := regexp.Compile(`import ([^java][a-zA-Z0-9]*\.[a-zA-Z0-9]*)`)
	if err != nil {
		return nil, err
	}
	// regex to find imports like org.springframework.boot
	regex2, err := regexp.Compile(`import ([^java][a-zA-Z0-9]*\.[a-zA-Z0-9]*\.[a-zA-Z0-9]*)`)
	if err != nil {
		return nil, err
	}
	// regex to find static imports like org.springframework (excluding standart java libraries)
	regex3, err := regexp.Compile(`import static ([^java][a-zA-Z0-9]*\.[a-zA-Z0-9]*)`)
	if err != nil {
		return nil, err
	}
	// regex to find static imports like org.springframework.boot
	regex4, err := regexp.Compile(`import static ([^java][a-zA-Z0-9]*\.[a-zA-Z0-9]*\.[a-zA-Z0-9]*)`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{regex1, regex2, regex3, regex4}), nil
}
