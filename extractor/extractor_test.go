package extractor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/codersrank-org/repo_info_extractor/extractor"
)

var _ = Describe("Extractor", func() {

	Context("RepoExtractor headless", func() {
		It("should the repo name with the owner name", func() {
			re := extractor.RepoExtractor{
				Headless: true,
			}
			Expect(re.GetRepoName("git@github.com:alimgiray/repo_info_extractor.git")).To(Equal("alimgiray/repo_info_extractor"))
			Expect(re.GetRepoName("https://github.com/alimgiray/repo_info_extractor.git")).To(Equal("alimgiray/repo_info_extractor"))
			Expect(re.GetRepoName("https://github.com/peti2001-test/second-project.git")).To(Equal("peti2001-test/second-project"))
		})
	})

	Context("RepoExtractor interactive", func() {
		It("should the repo name without the owner name", func() {
			re := extractor.RepoExtractor{
				Headless: false,
				RepoPath: "/some/path/alimgiray/repo_info_extractor",
			}
			Expect(re.GetRepoName("git@github.com:alimgiray/repo_info_extractor.git")).To(Equal("repo_info_extractor"))
			Expect(re.GetRepoName("https://github.com/alimgiray/repo_info_extractor.git")).To(Equal("repo_info_extractor"))
			Expect(re.GetRepoName("")).To(Equal("repo_info_extractor"))
		})
	})
})
