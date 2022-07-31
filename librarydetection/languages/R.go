package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection"
)

// NewRubyScriptAnalyzer constructor
func NewRAnalyzer() librarydetection.Analyzer {
	return &rAnalyzer{}
}

type rAnalyzer struct{}

func (a *rAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	// According to https://community.rstudio.com/t/analysis-package-naming-or-can-package-name-differ-from-rproj-name/9421
	// Package name can contain letters, numbers and coma

	// for this format library(shiny)
	requireRegex, err := regexp.Compile(`(?m)^\s*(?:require|library|package)\s*\(\s*["]*([\w\d,]+)"*\s*\)`)
	if err != nil {
		return nil, err
	}
	// for this format BiocManager123::repositories()
	requireRegex2, err := regexp.Compile(`(?m)^\s*([\w\d,]+)::`)
	if err != nil {
		return nil, err
	}
	// for this format library(package=quantmod)
	requireRegex3, err := regexp.Compile(`(?m)^\s*(?:require|library|package)\s*\(\s*package\s*=\s*([\w\d,]+)\s*\)`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{requireRegex, requireRegex2, requireRegex3}), nil
}
