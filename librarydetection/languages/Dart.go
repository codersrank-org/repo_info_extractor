package languages

import (
	"regexp"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection"
)

// NewDartScriptAnalyzer constructor
func NewDartScriptAnalyzer() librarydetection.Analyzer {
	return &dartScriptAnalyzer{}
}

type dartScriptAnalyzer struct{}

func (a *dartScriptAnalyzer) ExtractLibraries(contents string) ([]string, error) {
	// for imports like this: import 'dart:developer' as dev; OR import 'dart:developer';
	importRegex, err := regexp.Compile(`import '([:a-zA-Z0-9_-]+)'(?:|as)`)
	if err != nil {
		return nil, err
	}
	// for imports like this: import 'package:flutter/material.dart';
	importRegex2, err := regexp.Compile(`import '([:a-zA-Z0-9_-]+)/[\.a-zA-Z0-9_-]+'`)
	if err != nil {
		return nil, err
	}

	return executeRegexes(contents, []*regexp.Regexp{importRegex, importRegex2}), nil
}
