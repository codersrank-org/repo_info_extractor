package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection/languages"
)

var _ = Describe("PerlLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/perl.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"strict",
		"Benchmark",
		"Carp",
		"sigtrap",
		"Sub::Module",
		"Import::This",
		"utf8",
		"warnings",
		"Foo::Bar",
		"Module",
	}

	analyzer := languages.NewPerlAnalyzer()

	Describe("Extract Perl Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
