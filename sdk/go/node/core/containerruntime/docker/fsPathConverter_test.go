package docker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/util/iruntime"
)

var _ = Context("fsPathConverter", func() {
	Context("when runtime.GOOS == windows", func() {

		fakeRuntime := new(iruntime.Fake)
		fakeRuntime.GOOSReturns("windows")
		objectUnderTest := _fsPathConverter{
			runtime: fakeRuntime,
		}

		Context("when path contains drive letter", func() {
			It("should prepend a slash", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				actual := objectUnderTest.LocalToEngine("c:/DummyPath")

				/* assert */
				Expect(actual).To(Equal(expected))
			})
			It("should convert the drive letter to lowercase", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				actual := objectUnderTest.LocalToEngine("C:/DummyPath")

				/* assert */
				Expect(actual).To(Equal(expected))
			})
			It("should strip colon from drive", func() {
				/* arrange */
				expected := "/c/DummyPath"

				/* act */
				result := objectUnderTest.LocalToEngine("c:/DummyPath")

				/* assert */
				Expect(result).To(Equal(expected))
			})
			It("should replace backslashes with forward slashes", func() {
				/* arrange */
				pathWithMultipleBackslashes := `c\\dummy\path`
				expected := `c//dummy/path`

				/* act */
				result := objectUnderTest.LocalToEngine(pathWithMultipleBackslashes)

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
				actual := objectUnderTest.LocalToEngine(pathWithMultipleBackslashes)

				/* assert */
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
