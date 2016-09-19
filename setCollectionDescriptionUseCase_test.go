package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
  "errors"
  "reflect"
  "path"
)

var _ = Describe("_setCollectionDescriptionUseCase", func() {

  Context("Execute", func() {

    Context("when Filesystem.GetBytesOfFile returns an error", func() {

      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("GetBytesOfFileError")

        fakeFilesystem := new(FakeFilesystem)
        fakeFilesystem.GetBytesOfFileReturns(nil, expectedError)

        objectUnderTest := newSetCollectionDescriptionUseCase(
          fakeFilesystem,
          new(fakeFormat),
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetCollectionDescriptionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })

    })

    Context("when YamlFormat.From returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("FromError")

        fakeYamlFormat := new(fakeFormat)
        fakeYamlFormat.ToReturns(expectedError)

        objectUnderTest := newSetCollectionDescriptionUseCase(
          new(FakeFilesystem),
          fakeYamlFormat,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetCollectionDescriptionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    Context("when YamlFormat.To returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("ToError")

        fakeYamlFormat := new(fakeFormat)
        fakeYamlFormat.FromReturns(nil, expectedError)

        objectUnderTest := newSetCollectionDescriptionUseCase(
          new(FakeFilesystem),
          fakeYamlFormat,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetCollectionDescriptionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlFormat.From with expected collectionManifest", func() {

      /* arrange */
      expectedCollectionManifest := models.CollectionManifest{
        Manifest: models.Manifest{
          Name:"dummyName",
          Description:"dummyDescription",
          Version:"dummyVersion",
        },
      }

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedCollectionManifest))
        return
      }

      objectUnderTest := newSetCollectionDescriptionUseCase(
        new(FakeFilesystem),
        fakeYamlFormat,
      )

      /* act */
      objectUnderTest.Execute(
        models.SetCollectionDescriptionReq{Description:expectedCollectionManifest.Description},
      )

      /* assert */
      actualCollectionManifest := fakeYamlFormat.FromArgsForCall(0)
      Expect(actualCollectionManifest).To(Equal(&expectedCollectionManifest))

    })

    It("should call Filesystem.SaveFile with expected args", func() {

      /* arrange */
      providedPathToCollection := "/dummy/collection/path"
      expectedSaveFilePathArg := path.Join(providedPathToCollection, NameOfCollectionManifestFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := newSetCollectionDescriptionUseCase(
        fakeFilesystem,
        fakeYamlFormat,
      )

      /* act */
      objectUnderTest.Execute(
        models.SetCollectionDescriptionReq{PathToCollection:providedPathToCollection},
      )

      /* assert */
      actualSaveFilePathArg, actualSaveFileBytesArg := fakeFilesystem.SaveFileArgsForCall(0)
      Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
      Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

    })

  })

})
