package languages_test

import (
	. "github.com/onsi/ginkgo"
	"io/ioutil"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection/languages"
)

var _ = Describe("RustLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/rust.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"lib1",
		"lib2",
		"lib3",
		"lib4",
		"lib5",
		"lib6",
	}

	analyzer := languages.NewRustAnalyzer()

	Describe("Extract Rust Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
