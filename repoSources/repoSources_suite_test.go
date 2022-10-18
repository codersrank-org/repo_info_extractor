package repoSource_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRepoSources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RepoSources Suite")
}
