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
        new(fakeYamlCodec),
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
          new(fakeYamlCodec),
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.CreateCollectionReq{},
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

        objectUnderTest := newCreateCollectionUseCase(
          new(FakeFilesystem),
          fakeYamlCodec,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.CreateCollectionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlCodec.ToYaml with expected collectionManifest", func() {

      /* arrange */
      expectedCollectionManifest := models.CollectionManifest{
        Manifest:models.Manifest{
          Description:"DummyDescription",
          Name:"DummyName",
        },
      }

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedCollectionManifest))
        return
      }

      objectUnderTest := newCreateCollectionUseCase(
        new(FakeFilesystem),
        fakeYamlCodec,
      )

      /* act */
      objectUnderTest.Execute(
        models.CreateCollectionReq{
          Description:expectedCollectionManifest.Description,
          Name:expectedCollectionManifest.Name,
        },
      )

      /* assert */
      actualCollectionManifest := fakeYamlCodec.ToYamlArgsForCall(0)
      Expect(actualCollectionManifest).To(Equal(&expectedCollectionManifest))

    })

    It("should call Filesystem.SaveFile with expected args", func() {

      /* arrange */
      providedPath := "/dummy/op/path"
      expectedSaveFilePathArg := path.Join(providedPath, NameOfCollectionManifestFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.ToYamlReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := newCreateCollectionUseCase(
        fakeFilesystem,
        fakeYamlCodec,
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
