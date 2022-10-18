package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection/languages"
)

var _ = Describe("RubyLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/ruby.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"lib1",
		"lib2",
		"lib1/lib2",
		"lib-2/lib_1",
		"lib3",
		"lib4",
		"lib3/lib4",
		"lib4/lib3",
		"lib-13",
		"lib14",
		"lib15",
		"lib_16",
		"lib17",
		"lib18",
		"lib19",
	}

	analyzer := languages.NewRubyScriptAnalyzer()

	Describe("Extract Ruby Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
