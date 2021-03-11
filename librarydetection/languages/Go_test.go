package languages_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/codersrank-org/repo_info_extractor/librarydetection/languages"
)

var _ = Describe("GoLibraryDetection", func() {

	goFile := `
		import (
			"context"
			"log"
			"github.com/codersrank-org/repo_info_extractor"
		)
	`
	libraries := []string{"context", "log", "github.com/codersrank-org/repo_info_extractor"}
	goAnalyzer := languages.NewGoAnalyzer()

	Describe("ExtractLibraries", func() {
		It("Should be able to extract libraries", func() {
			libs := goAnalyzer.ExtractLibraries(goFile)
			Expect(len(libs)).To(Equal(3))
			Expect(libs[0]).To(Equal(libraries[0]))
			Expect(libs[1]).To(Equal(libraries[1]))
			Expect(libs[2]).To(Equal(libraries[2]))
		})
	})

})
