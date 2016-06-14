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
          new(fakeYamlCodec),
        )

        /* act */
        _, actualError := objectUnderTest.Construct("/dummy/path")

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })

    })

    Context("when YamlCodec.FromYaml returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("FromYamlError")

        fakeYamlCodec := new(fakeYamlCodec)
        fakeYamlCodec.FromYamlReturns(expectedError)

        objectUnderTest := newCollectionViewFactory(
          new(FakeFilesystem),
          new(fakeOpViewFactory),
          fakeYamlCodec,
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
          new(fakeYamlCodec),
        )

        /* act */
        _, actualError := objectUnderTest.Construct("/dummy/path")

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlCodec.FromYaml with expected bytes", func() {

      /* arrange */
      expectedBytes := []byte{0, 8, 10}

      fakeFilesystem := new(FakeFilesystem)
      fakeFilesystem.GetBytesOfFileReturns(expectedBytes, nil)

      fakeYamlCodec := new(fakeYamlCodec)

      objectUnderTest := newCollectionViewFactory(
        fakeFilesystem,
        new(fakeOpViewFactory),
        fakeYamlCodec,
      )

      /* act */
      objectUnderTest.Construct("/dummy/path")

      /* assert */
      actualBytes, _ := fakeYamlCodec.FromYamlArgsForCall(0)
      Expect(actualBytes).To(Equal(expectedBytes))

    })

  })

})
