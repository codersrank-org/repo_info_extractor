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
		It("should detect build files ", func() {
			// Act
			l1 := a.Detect("/home/something/Makefile", []byte{})
			l2 := a.Detect("/home/something/Dockerfile", []byte{})
			l3 := a.Detect("/home/something/Jenkinsfile", []byte{})
			l4 := a.Detect("/home/something/Rakefile", []byte{})
			l5 := a.Detect("/home/something/CMakeLists.txt", []byte{})

			// Assert
			Expect(l1).To(Equal("Makefile"))
			Expect(l2).To(Equal("Dockerfile"))
			Expect(l3).To(Equal("Jenkins"))
			Expect(l4).To(Equal("Ruby"))
			Expect(l5).To(Equal("CMake"))
		})
	})

	Context("Detect language by file name", func() {
		It("should detect SQL and PLpgSQL ", func() {
			// Act
			l1 := a.Detect("/home/something/get_pg_users.sql", []byte(`SELECT usename FROM pg_catalog.pg_user;`))
			l2 := a.Detect("/home/something/create_pg_user.sql", []byte(`--
CREATE OR REPLACE FUNCTION __tmp_create_user()
  RETURNS VOID
  LANGUAGE plpgsql
AS
$$
BEGIN
  IF NOT EXISTS (
      SELECT
      FROM   pg_catalog.pg_user
      WHERE  usename = 'new_user') THEN
    CREATE USER 'new_user';
  END IF;
END;
$$;`))

			// Assert
			Expect(l1).To(Equal("SQL"))
			Expect(l2).To(Equal("PLpgSQL"))
		})
	})
})
