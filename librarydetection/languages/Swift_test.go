package languages_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"

	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
)

var _ = Describe("SwiftLibraryDetection", func() {
	fixture, err := ioutil.ReadFile("./fixtures/swift.fixture")
	if err != nil {
		panic(err)
	}

	expectedLibraries := []string{
		"Cocoa",
		"os.log",
		"StatsKit",
		"ModuleKit",
		"CPU",
		"Memory",
		"Disk",
		"Net",
		"Battery",
		"Sensors",
		"GPU",
		"Fans",
		"Pentathlon.swim",
		"test.test",
		"CoreServices.DictionaryServices",
	}

	analyzer := languages.NewSwiftAnalyzer()

	Describe("Extract Swift Libraries", func() {
		It("Should be able to extract libraries", func() {
			libs, err := analyzer.ExtractLibraries(string(fixture))
			if err != nil {
				panic(err)
			}
			assertSameUnordered(libs, expectedLibraries)
		})
	})
})
