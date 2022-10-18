package repoSource

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Local", func() {
	Describe("GetRepos", func() {
		It("should set the name to the custom name", func() {
			// Arrange
			localSource := NewDirectoryPath("/path/to/test/repo", "newRepoName")

			// Act
			repos := localSource.GetRepos()

			// Assert
			Expect(len(repos)).To(Equal(1))
			Expect(repos[0].FullName).To(Equal("newRepoName"))
			Expect(repos[0].GetSafeFullName()).To(Equal("newRepoName"))
		})
		It("should set the name based on the path", func() {
			// Arrange
			localSource := NewDirectoryPath("/path/to/test/repo_name", "")

			// Act
			repos := localSource.GetRepos()

			// Assert
			Expect(len(repos)).To(Equal(1))
			Expect(repos[0].FullName).To(Equal("repo_name"))
			Expect(repos[0].GetSafeFullName()).To(Equal("repo_name"))
		})
	})
})
