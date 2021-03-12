package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
)

var _ = Describe("JavaScriptLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/javascript.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"lib1",
		"lib2",
		"lib3",
		"lib4",
	}

	analyzer := languages.NewJavaScriptAnalyzer()

	Describe("Extract JavaScript Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs := analyzer.ExtractLibraries(string(fixture))
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
