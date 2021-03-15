package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
)

var _ = Describe("HTMLLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/html.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"style.css",
		"bootstrap.css",
		"jquery.min.js",
		"hello.js",
	}

	analyzer := languages.NewHTMLAnalyzer()

	Describe("Extract HTML Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
