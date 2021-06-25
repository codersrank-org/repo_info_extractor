package languagedetection

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Analyzer", func() {
	a := NewLanguageAnalyzer()
	Context("Detect language by extension", func() {
		It("should detect PHP ", func() {
			// Act
			l1 := a.Detect("/home/something/index.php", []byte{})
			l2 := a.Detect("/home/something/index.Php", []byte{})
			l3 := a.Detect("/home/something/index.razor", []byte{})

			// Assert
			Expect(l1).To(Equal("PHP"))
			Expect(l2).To(Equal("PHP"))
			Expect(l3).To(Equal("Blazor"))
		})
	})

	Context("Detect language by file name", func() {
		It("should detect PHP ", func() {
			// Act
			l1 := a.Detect("/home/something/Makefile", []byte{})
			l2 := a.Detect("/home/something/Dockerfile", []byte{})

			// Assert
			Expect(l1).To(Equal("Makefile"))
			Expect(l2).To(Equal("Dockerfile"))
		})
	})
})
