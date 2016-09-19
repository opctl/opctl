package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
  "errors"
  "reflect"
  "path"
)

var _ = Describe("_setOpDescriptionUseCase", func() {

  Context("Execute", func() {

    Context("when Filesystem.GetBytesOfFile returns an error", func() {

      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("GetBytesOfFileError")

        fakeFilesystem := new(FakeFilesystem)
        fakeFilesystem.GetBytesOfFileReturns(nil, expectedError)

        objectUnderTest := newSetOpDescriptionUseCase(
          fakeFilesystem,
          new(fakeFormat),
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetOpDescriptionReq{},
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

        objectUnderTest := newSetOpDescriptionUseCase(
          new(FakeFilesystem),
          fakeYamlFormat,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetOpDescriptionReq{},
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

        objectUnderTest := newSetOpDescriptionUseCase(
          new(FakeFilesystem),
          fakeYamlFormat,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetOpDescriptionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlFormat.From with expected opManifest", func() {

      /* arrange */
      expectedOpManifest := models.OpManifest{
        Manifest: models.Manifest{
          Name:"dummyName",
          Description:"dummyDescription",
          Version:"dummyVersion",
        },
      }

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedOpManifest))
        return
      }

      objectUnderTest := newSetOpDescriptionUseCase(
        new(FakeFilesystem),
        fakeYamlFormat,
      )

      /* act */
      objectUnderTest.Execute(
        models.SetOpDescriptionReq{Description:expectedOpManifest.Description},
      )

      /* assert */
      actualOpManifest := fakeYamlFormat.FromArgsForCall(0)
      Expect(actualOpManifest).To(Equal(&expectedOpManifest))

    })

    It("should call Filesystem.SaveFile with expected args", func() {

      /* arrange */
      providedPathToOp := "/dummy/op/path"
      expectedSaveFilePathArg := path.Join(providedPathToOp, NameOfOpManifestFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := newSetOpDescriptionUseCase(
        fakeFilesystem,
        fakeYamlFormat,
      )

      /* act */
      objectUnderTest.Execute(
        models.SetOpDescriptionReq{PathToOp:providedPathToOp},
      )

      /* assert */
      actualSaveFilePathArg, actualSaveFileBytesArg := fakeFilesystem.SaveFileArgsForCall(0)
      Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
      Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

    })

  })

})
