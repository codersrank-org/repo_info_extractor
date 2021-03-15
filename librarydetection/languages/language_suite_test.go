package languages_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAccount(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Suite")
}

// assertSameUnordered asserts that two slices of strings are the same, regardless of order
func assertSameUnordered(result, expected []string) {
	Expect(len(result)).Should(Equal(len(expected)))
	for _, v := range result {
		Expect(v).Should(BeElementOf(expected))
	}
}
