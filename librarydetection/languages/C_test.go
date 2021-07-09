package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection/languages"
)

var _ = Describe("CLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/c.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"assert.h",
		"complex.h",
		"ctype.h",
		"stdio.h",
		"float.h",
		"string.h",
		"hey/sup/iomanip.h",
		"Hey/ssup/math.h",
		"hello/how/stdlib.h",
		"great/wchar.h",
		"stdbool.h",
		"stdint.h",
		"hello12",
		"WsSup34",
		"heyYo3-lol",
	}

	analyzer := languages.NewCAnalyzer()

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
