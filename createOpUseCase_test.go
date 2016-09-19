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
        new(fakeFormat),
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
          new(fakeFormat),
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.CreateOpReq{},
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
        fakeYamlFormat.FromReturns(nil, expectedError)

        objectUnderTest := newCreateOpUseCase(
          new(FakeFilesystem),
          fakeYamlFormat,
        )

        /* act */
        actualError := objectUnderTest.Execute(
          models.CreateOpReq{},
        )

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlFormat.From with expected opManifest", func() {

      /* arrange */
      expectedOpManifest := models.OpManifest{
        Manifest:models.Manifest{
          Description:"DummyDescription",
          Name:"DummyName",
        },
      }

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedOpManifest))
        return
      }

      objectUnderTest := newCreateOpUseCase(
        new(FakeFilesystem),
        fakeYamlFormat,
      )

      /* act */
      objectUnderTest.Execute(
        models.CreateOpReq{
          Description:expectedOpManifest.Description,
          Name:expectedOpManifest.Name,
        },
      )

      /* assert */
      actualOpManifest := fakeYamlFormat.FromArgsForCall(0)
      Expect(actualOpManifest).To(Equal(&expectedOpManifest))

    })

    It("should call Filesystem.SaveFile with expected args", func() {

      /* arrange */
      providedPath := "/dummy/op/path"
      expectedSaveFilePathArg := path.Join(providedPath, NameOfOpManifestFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := newCreateOpUseCase(
        fakeFilesystem,
        fakeYamlFormat,
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
