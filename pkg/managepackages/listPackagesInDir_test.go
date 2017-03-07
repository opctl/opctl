package managepackages

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
)

var _ = Describe("listPackagesInDir", func() {

	Context("ListPackagesInDir", func() {

		Context("when FileSystem.ListChildFileInfosOfDir returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("ListChildFileInfosOfDirError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.ListChildFileInfosOfDirReturns(nil, expectedError)

				objectUnderTest := managePackages{
					fileSystem:         fakeFileSystem,
					packageViewFactory: new(fakePackageViewFactory),
					yaml:               new(format.Fake),
				}

				/* act */
				_, actualError := objectUnderTest.ListPackagesInDir("/dummy/path")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})
	})

})
