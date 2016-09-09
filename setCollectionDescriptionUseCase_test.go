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
          new(fakeYamlCodec),
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetCollectionDescriptionReq{},
        )

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

        objectUnderTest := newSetCollectionDescriptionUseCase(
          new(FakeFilesystem),
          fakeYamlCodec,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetCollectionDescriptionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    Context("when YamlCodec.ToYaml returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("ToYamlError")

        fakeYamlCodec := new(fakeYamlCodec)
        fakeYamlCodec.ToYamlReturns(nil, expectedError)

        objectUnderTest := newSetCollectionDescriptionUseCase(
          new(FakeFilesystem),
          fakeYamlCodec,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetCollectionDescriptionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlCodec.ToYaml with expected collectionManifest", func() {

      /* arrange */
      expectedCollectionManifest := models.CollectionManifest{
        Manifest: models.Manifest{
          Name:"dummyName",
          Description:"dummyDescription",
          Version:"dummyVersion",
        },
      }

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedCollectionManifest))
        return
      }

      objectUnderTest := newSetCollectionDescriptionUseCase(
        new(FakeFilesystem),
        fakeYamlCodec,
      )

      /* act */
      objectUnderTest.Execute(
        models.SetCollectionDescriptionReq{Description:expectedCollectionManifest.Description},
      )

      /* assert */
      actualCollectionManifest := fakeYamlCodec.ToYamlArgsForCall(0)
      Expect(actualCollectionManifest).To(Equal(&expectedCollectionManifest))

    })

    It("should call Filesystem.SaveFile with expected args", func() {

      /* arrange */
      providedPathToCollection := "/dummy/collection/path"
      expectedSaveFilePathArg := path.Join(providedPathToCollection, NameOfCollectionManifestFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.ToYamlReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := newSetCollectionDescriptionUseCase(
        fakeFilesystem,
        fakeYamlCodec,
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
