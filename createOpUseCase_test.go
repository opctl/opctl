package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
  "errors"
  "reflect"
  "path"
)

var _ = Describe("_createOpUseCase", func() {

  Context("Execute", func() {

    It("should call Filesystem.AddDir with expected args", func() {

      /* arrange */

      providedCreateOpReq := models.CreateOpReq{Path:"/dummy/path"}

      fakeFilesystem := new(FakeFilesystem)

      objectUnderTest := newCreateOpUseCase(
        fakeFilesystem,
        new(fakeYamlCodec),
      )

      /* act */
      objectUnderTest.Execute(
        providedCreateOpReq,
      )

      /* assert */
      Expect(fakeFilesystem.AddDirArgsForCall(0)).To(Equal(providedCreateOpReq.Path))

    })

    Context("when Filesystem.AddDir returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("AddDirError")

        fakeFilesystem := new(FakeFilesystem)
        fakeFilesystem.AddDirReturns(expectedError)

        objectUnderTest := newCreateOpUseCase(
          fakeFilesystem,
          new(fakeYamlCodec),
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.CreateOpReq{},
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

        objectUnderTest := newCreateOpUseCase(
          new(FakeFilesystem),
          fakeYamlCodec,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.CreateOpReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlCodec.ToYaml with expected opBundleManifest", func() {

      /* arrange */
      expectedOpBundleManifest := models.OpBundleManifest{
        BundleManifest:models.BundleManifest{
          Description:"DummyDescription",
          Name:"DummyName",
        },
      }

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedOpBundleManifest))
        return
      }

      objectUnderTest := newCreateOpUseCase(
        new(FakeFilesystem),
        fakeYamlCodec,
      )

      /* act */
      objectUnderTest.Execute(
        models.CreateOpReq{
          Description:expectedOpBundleManifest.Description,
          Name:expectedOpBundleManifest.Name,
        },
      )

      /* assert */
      actualOpBundleManifest := fakeYamlCodec.ToYamlArgsForCall(0)
      Expect(actualOpBundleManifest).To(Equal(&expectedOpBundleManifest))

    })

    It("should call Filesystem.SaveFile with expected args", func() {

      /* arrange */
      providedPath := "/dummy/op/path"
      expectedSaveFilePathArg := path.Join(providedPath, NameOfOpBundleManifest)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.ToYamlReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := newCreateOpUseCase(
        fakeFilesystem,
        fakeYamlCodec,
      )

      /* act */
      objectUnderTest.Execute(
        models.CreateOpReq{Path:providedPath},
      )

      /* assert */
      actualSaveFilePathArg, actualSaveFileBytesArg := fakeFilesystem.SaveFileArgsForCall(0)
      Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
      Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

    })

  })

})
