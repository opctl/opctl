package pathnormalizer

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("pathNormalizer", func() {
  Context("Normalize()", func() {
    Describe("when path contains drive letter", func() {
      It("should prepend a slash", func() {
        /* arrange */
        expected := "/c/DummyPath"

        objectUnderTest := NewPathNormalizer()

        /* act */
        actual := objectUnderTest.Normalize("c:/DummyPath")

        /* assert */
        Expect(actual).To(Equal(expected))
      })
      It("should convert the drive letter to lowercase", func() {
        /* arrange */
        expected := "/c/DummyPath"

        objectUnderTest := NewPathNormalizer()

        /* act */
        actual := objectUnderTest.Normalize("C:/DummyPath")

        /* assert */
        Expect(actual).To(Equal(expected))
      })
      It("should strip colon from drive", func() {
        /* arrange */
        expected := "/c/DummyPath"

        objectUnderTest := NewPathNormalizer()

        /* act */
        result := objectUnderTest.Normalize("c:/DummyPath")

        /* assert */
        Expect(result).To(Equal(expected))
      })
      It("should replace single backslashes with single forward slashes", func() {
        /* arrange */
        pathWithMultipleBackslashes := `c\dummy\path`
        expected := `c/dummy/path`

        objectUnderTest := NewPathNormalizer()

        /* act */
        result := objectUnderTest.Normalize(pathWithMultipleBackslashes)

        /* assert */
        Expect(result).To(Equal(expected))
      })
      It("should replace double backslashes with single forward slashes", func() {
        /* arrange */
        pathWithMultipleBackslashes := `c\\dummy\\path`
        expected := `c/dummy/path`

        objectUnderTest := NewPathNormalizer()

        /* act */
        actual := objectUnderTest.Normalize(pathWithMultipleBackslashes)

        /* assert */
        Expect(actual).To(Equal(expected))
      })
    })
    Describe("when path doesn't contain a drive letter", func() {
      It("should replace single backslashes with single forward slashes", func() {
        /* arrange */
        pathWithMultipleBackslashes := `\dummy\path`
        expected := `/dummy/path`

        objectUnderTest := NewPathNormalizer()

        /* act */
        actual := objectUnderTest.Normalize(pathWithMultipleBackslashes)

        /* assert */
        Expect(actual).To(Equal(expected))
      })
      It("should replace double backslashes with single forward slashes", func() {
        /* arrange */
        pathWithMultipleBackslashes := `\\dummy\\path`
        expected := `/dummy/path`

        objectUnderTest := NewPathNormalizer()

        /* act */
        actual := objectUnderTest.Normalize(pathWithMultipleBackslashes)

        /* assert */
        Expect(actual).To(Equal(expected))
      })
    })
  })
})
