package extractor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/codersrank-org/repo_info_extractor/v2/extractor"
)

var _ = Describe("GetRepoName", func() {

	Context("RepoExtractor headless", func() {
		It("should get the repo name with the owner name", func() {
			re := extractor.RepoExtractor{
				Headless: true,
			}
			Expect(re.GetRepoName("git@github.com:alimgiray/repo_info_extractor.git")).To(Equal("alimgiray/repo_info_extractor"))
			Expect(re.GetRepoName("https://github.com/alimgiray/repo_info_extractor.git")).To(Equal("alimgiray/repo_info_extractor"))
			Expect(re.GetRepoName("https://github.com/peti2001-test/second-project.git")).To(Equal("peti2001-test/second-project"))
			Expect(re.GetRepoName("ssh://user@host:port/group/repoName.git")).To(Equal("group/repoName"))
		})
	})

	Context("RepoExtractor interactive", func() {
		It("should get the repo name without the owner name", func() {
			re := extractor.RepoExtractor{
				Headless: false,
				RepoPath: "/some/path/alimgiray/repo_info_extractor",
			}
			Expect(re.GetRepoName("git@github.com:alimgiray/repo_info_extractor.git")).To(Equal("repo_info_extractor"))
			Expect(re.GetRepoName("https://github.com/alimgiray/repo_info_extractor.git")).To(Equal("repo_info_extractor"))
			Expect(re.GetRepoName("")).To(Equal("repo_info_extractor"))
			Expect(re.GetRepoName("ssh://user@host:port/group/repoName.git")).To(Equal("repoName"))
		})
	})
})
