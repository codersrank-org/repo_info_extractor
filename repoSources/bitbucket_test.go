package repoSource

import (
	"io/ioutil"
	"os"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getResponseFromFile(filePath string) []byte {
	responseFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer responseFile.Close()

	byteValue, err := ioutil.ReadAll(responseFile)
	if err != nil {
		panic(err)
	}
	return byteValue
}

var _ = Describe("Bitbucket", func() {
	Describe("Getting repositories", func() {
		It("should get repositories of the user", func() {
			// Arrange
			httpmock.Activate()
			httpmock.RegisterResponder("GET", "https://api.bitbucket.org/2.0/repositories?q=is_private+%3D+false&role=contributor", httpmock.NewStringResponder(200, string(getResponseFromFile("../test_fixtures/repoSources/bitbucket/public.json"))))
			provider := NewBitbucketProvider("test_username", "test_token", "public", "")

			// Act
			repos := provider.GetRepos()

			// Assert
			Expect(len(repos)).To(Equal(10))
			Expect(repos[0].FullName).To(Equal("opensymphony/xwork"))
			Expect(repos[0].Name).To(Equal("xwork"))
			Expect(repos[0].ID).To(Equal("{3f630668-75f1-4903-ae5e-8ea37437e09e}"))
			Expect(repos[0].CloneURL).To(Equal("https://bitbucket.org/opensymphony/xwork.git"))
			httpmock.DeactivateAndReset()
		})
	})

})
