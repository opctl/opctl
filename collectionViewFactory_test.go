package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "errors"
)

var _ = Describe("_collectionViewFactory", func() {

  Context("Construct", func() {

    Context("when Filesystem.GetBytesOfFile returns an error", func() {

      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("GetBytesOfFileError")

        fakeFilesystem := new(FakeFilesystem)
        fakeFilesystem.GetBytesOfFileReturns(nil, expectedError)

        objectUnderTest := newCollectionViewFactory(
          fakeFilesystem,
          new(fakeOpViewFactory),
          new(fakeFormat),
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

        fakeYamlFormat := new(fakeFormat)
        fakeYamlFormat.ToReturns(expectedError)

        objectUnderTest := newCollectionViewFactory(
          new(FakeFilesystem),
          new(fakeOpViewFactory),
          fakeYamlFormat,
        )

        /* act */
        _, actualError := objectUnderTest.Construct("/dummy/path")

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    Context("when Filesystem.ListChildFileInfosOfDir returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("ListChildFileInfosOfDirError")

        fakeFilesystem := new(FakeFilesystem)
        fakeFilesystem.ListChildFileInfosOfDirReturns(nil, expectedError)

        objectUnderTest := newCollectionViewFactory(
          fakeFilesystem,
          new(fakeOpViewFactory),
          new(fakeFormat),
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

      fakeFilesystem := new(FakeFilesystem)
      fakeFilesystem.GetBytesOfFileReturns(expectedBytes, nil)

      fakeYamlFormat := new(fakeFormat)

      objectUnderTest := newCollectionViewFactory(
        fakeFilesystem,
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
