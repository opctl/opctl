package docker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/vruntime"
)

var _ = Context("enginePath", func() {
	Context("when runtime.GOOS == windows", func() {

		fakeRuntime := new(vruntime.Fake)
		fakeRuntime.GOOSReturns("windows")
		objectUnderTest := _containerProvider{
			runtime: fakeRuntime,
		}

		Context("when path contains drive letter", func() {
			It("should prepend a slash", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				actual := objectUnderTest.enginePath("c:/DummyPath")

				/* assert */
				Expect(actual).To(Equal(expected))
			})
			It("should convert the drive letter to lowercase", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				actual := objectUnderTest.enginePath("C:/DummyPath")

				/* assert */
				Expect(actual).To(Equal(expected))
			})
			It("should strip colon from drive", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				result := objectUnderTest.enginePath("c:/DummyPath")

				/* assert */
				Expect(result).To(Equal(expected))
			})
			It("should replace backslashes with forward slashes", func() {
				/* arrange */
				pathWithMultipleBackslashes := `c\\dummy\path`
				expected := `c//dummy/path`

				/* act */
				result := objectUnderTest.enginePath(pathWithMultipleBackslashes)

				/* assert */
				Expect(result).To(Equal(expected))
			})
		})
		Context("when path doesn't contain a drive letter", func() {
			It("should replace backslashes with forward slashes", func() {
				/* arrange */
				pathWithMultipleBackslashes := `\\dummy\path`
				expected := `//dummy/path`

				/* act */
				actual := objectUnderTest.enginePath(pathWithMultipleBackslashes)

				/* assert */
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
