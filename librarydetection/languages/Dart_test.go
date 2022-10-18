package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection/languages"
)

var _ = Describe("DartLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/dart.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"dart:async",
		"dart:developer",
		"package:flutter",
	}

	analyzer := languages.NewDartScriptAnalyzer()

	Describe("Extract Dart Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
