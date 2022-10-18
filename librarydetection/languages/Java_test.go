package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/v2/librarydetection/languages"
)

var _ = Describe("JavaLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/java.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"org.bytedeco",
		"org.bytedeco.javacv",
		"org.bytedeco",
		"org.bytedeco.opencv",
	}

	analyzer := languages.NewJavaAnalyzer()

	Describe("Extract Java Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
