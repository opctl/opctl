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
          new(fakeYamlCodec),
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetOpDescriptionReq{},
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

        objectUnderTest := newSetOpDescriptionUseCase(
          new(FakeFilesystem),
          fakeYamlCodec,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.SetOpDescriptionReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlCodec.ToYaml with expected opFile", func() {

      /* arrange */
      expectedOpFile := models.OpFile{
        Name:"DummyName",
        Description:"DummyDescription",
      }

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedOpFile))
        return
      }

      objectUnderTest := newSetOpDescriptionUseCase(
        new(FakeFilesystem),
        fakeYamlCodec,
      )

      /* act */
      objectUnderTest.Execute(
        models.SetOpDescriptionReq{Description:expectedOpFile.Description},
      )

      /* assert */
      actualOpFile := fakeYamlCodec.ToYamlArgsForCall(0)
      Expect(actualOpFile).To(Equal(&expectedOpFile))

    })

    It("should call Filesystem.SaveFile with expected args", func() {

      /* arrange */
      providedPathToOp := "/dummy/op/path"
      expectedSaveFilePathArg := path.Join(providedPathToOp, NameOfOpFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.ToYamlReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := newSetOpDescriptionUseCase(
        fakeFilesystem,
        fakeYamlCodec,
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
