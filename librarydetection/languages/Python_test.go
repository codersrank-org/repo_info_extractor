package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
)

var _ = Describe("PythonLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/python.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"lib1.lib2",
		"lib3",
		"lib4",
	}

	analyzer := languages.NewPythonScriptAnalyzer()

	Describe("Extract Python Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
