package languagedetection_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLanguagedetection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Languagedetection Suite")
}
