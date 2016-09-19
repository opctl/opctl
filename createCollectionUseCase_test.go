package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
  "errors"
  "reflect"
  "path"
)

var _ = Describe("_createCollectionUseCase", func() {

  Context("Execute", func() {

    It("should call Filesystem.AddDir with expected args", func() {

      /* arrange */

      providedCreateCollectionReq := models.CreateCollectionReq{Path:"/dummy/path"}

      fakeFilesystem := new(FakeFilesystem)

      objectUnderTest := newCreateCollectionUseCase(
        fakeFilesystem,
        new(fakeFormat),
      )

      /* act */
      objectUnderTest.Execute(
        providedCreateCollectionReq,
      )

      /* assert */
      Expect(fakeFilesystem.AddDirArgsForCall(0)).To(Equal(providedCreateCollectionReq.Path))

    })

    Context("when Filesystem.AddDir returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("AddDirError")

        fakeFilesystem := new(FakeFilesystem)
        fakeFilesystem.AddDirReturns(expectedError)

        objectUnderTest := newCreateCollectionUseCase(
          fakeFilesystem,
          new(fakeFormat),
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.CreateCollectionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    Context("when YamlFormat.From returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("ToError")

        fakeYamlFormat := new(fakeFormat)
        fakeYamlFormat.FromReturns(nil, expectedError)

        objectUnderTest := newCreateCollectionUseCase(
          new(FakeFilesystem),
          fakeYamlFormat,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.CreateCollectionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlFormat.From with expected collectionManifest", func() {

      /* arrange */
      expectedCollectionManifest := models.CollectionManifest{
        Manifest:models.Manifest{
          Description:"DummyDescription",
          Name:"DummyName",
        },
      }

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedCollectionManifest))
        return
      }

      objectUnderTest := newCreateCollectionUseCase(
        new(FakeFilesystem),
        fakeYamlFormat,
      )

      /* act */
      objectUnderTest.Execute(
        models.CreateCollectionReq{
          Description:expectedCollectionManifest.Description,
          Name:expectedCollectionManifest.Name,
        },
      )

      /* assert */
      actualCollectionManifest := fakeYamlFormat.FromArgsForCall(0)
      Expect(actualCollectionManifest).To(Equal(&expectedCollectionManifest))

    })

    It("should call Filesystem.SaveFile with expected args", func() {

      /* arrange */
      providedPath := "/dummy/op/path"
      expectedSaveFilePathArg := path.Join(providedPath, NameOfCollectionManifestFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := newCreateCollectionUseCase(
        fakeFilesystem,
        fakeYamlFormat,
      )

      /* act */
      objectUnderTest.Execute(
        models.CreateCollectionReq{Path:providedPath},
      )

      /* assert */
      actualSaveFilePathArg, actualSaveFileBytesArg := fakeFilesystem.SaveFileArgsForCall(0)
      Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
      Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

    })

  })

})
