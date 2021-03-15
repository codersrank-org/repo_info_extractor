package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
)

var _ = Describe("CLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/kotlin.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"syntax.lexer.Token",
		"syntax.tree",
	}

	analyzer := languages.NewKotlinAnalyzer()

	Describe("Extract C Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
