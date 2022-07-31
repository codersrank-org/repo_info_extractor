package languages_test

import (
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection/languages"
)

var _ = Describe("RLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/r.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"ggplot",
		"quantmod",
		"tibble",
		"tidyr",
		"shiny",
		"bslib",
		"DT",
		"WGCNA",
		"impute",
		"BiocManager123",
	}

	analyzer := languages.NewRAnalyzer()

	Describe("Extract R Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			for _, l := range libs {
				fmt.Println(l)
			}
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
