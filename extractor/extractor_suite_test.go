package extractor_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExtractor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extractor Suite")
}
