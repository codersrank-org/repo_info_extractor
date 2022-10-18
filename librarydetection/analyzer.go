package librarydetection

import (
	"fmt"
)

// Analyzer is an interface for extracting various features from files
// Language specific implementations are at ./languages folder
type Analyzer interface {
	ExtractLibraries(contents string) ([]string, error)
}

// Analyzers is the map for all analyzers
// like "Go" has "GoAnalyzer", "Python" has "PythonAnalyzer" and so on.
type Analyzers map[string]Analyzer

var analyzers = Analyzers{}

// GetAnalyzer returns given analyzer for that language
func GetAnalyzer(language string) (Analyzer, error) {
	analyzer := analyzers[language]
	if analyzer == nil {
		return nil, fmt.Errorf("no analyzer for %s exists", language)
	}
	return analyzer, nil
}

// AddAnalyzer allows users to add new analyzers
func AddAnalyzer(language string, analyzer Analyzer) {
	analyzers[language] = analyzer
}
