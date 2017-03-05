package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
)

var _ = Describe("_collectionViewFactory", func() {

	Context("Construct", func() {

		Context("when FileSystem.GetBytesOfFile returns an error", func() {

			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("GetBytesOfFileError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.GetBytesOfFileReturns(nil, expectedError)

				objectUnderTest := newCollectionViewFactory(
					fakeFileSystem,
					new(fakeOpViewFactory),
					new(format.Fake),
				)

				/* act */
				_, actualError := objectUnderTest.Construct("/dummy/path")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})

		})

		Context("when YamlFormat.To returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("FromError")

				fakeYamlFormat := new(format.Fake)
				fakeYamlFormat.ToReturns(expectedError)

				objectUnderTest := newCollectionViewFactory(
					new(fs.Fake),
					new(fakeOpViewFactory),
					fakeYamlFormat,
				)

				/* act */
				_, actualError := objectUnderTest.Construct("/dummy/path")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		Context("when FileSystem.ListChildFileInfosOfDir returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("ListChildFileInfosOfDirError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.ListChildFileInfosOfDirReturns(nil, expectedError)

				objectUnderTest := newCollectionViewFactory(
					fakeFileSystem,
					new(fakeOpViewFactory),
					new(format.Fake),
				)

				/* act */
				_, actualError := objectUnderTest.Construct("/dummy/path")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		It("should call YamlFormat.To with expected bytes", func() {

			/* arrange */
			expectedBytes := []byte{0, 8, 10}

			fakeFileSystem := new(fs.Fake)
			fakeFileSystem.GetBytesOfFileReturns(expectedBytes, nil)

			fakeYamlFormat := new(format.Fake)

			objectUnderTest := newCollectionViewFactory(
				fakeFileSystem,
				new(fakeOpViewFactory),
				fakeYamlFormat,
			)

			/* act */
			objectUnderTest.Construct("/dummy/path")

			/* assert */
			actualBytes, _ := fakeYamlFormat.ToArgsForCall(0)
			Expect(actualBytes).To(Equal(expectedBytes))

		})

	})

})
