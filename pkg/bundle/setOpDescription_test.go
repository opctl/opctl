package bundle

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "errors"
  "reflect"
  "path"
  "github.com/opspec-io/sdk-golang/util/fs"
  "github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("_setOpDescription", func() {

  Context("Execute", func() {

    Context("when FileSystem.GetBytesOfFile returns an error", func() {

      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("GetBytesOfFileError")

        fakeFileSystem := new(fs.FakeFileSystem)
        fakeFileSystem.GetBytesOfFileReturns(nil, expectedError)

        objectUnderTest := &_bundle{
          fileSystem: fakeFileSystem,
          yaml: format.NewYamlFormat(),
        }

        /* act */
        actualError := objectUnderTest.SetOpDescription(
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

        fakeYamlFormat := new(format.FakeFormat)
        fakeYamlFormat.ToReturns(expectedError)

        objectUnderTest := &_bundle{
          fileSystem: new(fs.FakeFileSystem),
          yaml: fakeYamlFormat,
        }

        /* act */
        actualError := objectUnderTest.SetOpDescription(
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

        fakeYamlFormat := new(format.FakeFormat)
        fakeYamlFormat.FromReturns(nil, expectedError)

        objectUnderTest := &_bundle{
          fileSystem: new(fs.FakeFileSystem),
          yaml: fakeYamlFormat,
        }

        /* act */
        actualError := objectUnderTest.SetOpDescription(
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

      fakeYamlFormat := new(format.FakeFormat)
      fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedOpManifest))
        return
      }

      objectUnderTest := &_bundle{
        fileSystem: new(fs.FakeFileSystem),
        yaml: fakeYamlFormat,
      }

      /* act */
      objectUnderTest.SetOpDescription(
        models.SetOpDescriptionReq{Description:expectedOpManifest.Description},
      )

      /* assert */
      actualOpManifest := fakeYamlFormat.FromArgsForCall(0)
      Expect(actualOpManifest).To(Equal(&expectedOpManifest))

    })

    It("should call FileSystem.SaveFile with expected args", func() {

      /* arrange */
      providedPathToOp := "/dummy/op/path"
      expectedSaveFilePathArg := path.Join(providedPathToOp, NameOfOpManifestFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFileSystem := new(fs.FakeFileSystem)

      fakeYamlFormat := new(format.FakeFormat)
      fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := &_bundle{
        fileSystem: fakeFileSystem,
        yaml: fakeYamlFormat,
      }

      /* act */
      objectUnderTest.SetOpDescription(
        models.SetOpDescriptionReq{PathToOp:providedPathToOp},
      )

      /* assert */
      actualSaveFilePathArg, actualSaveFileBytesArg := fakeFileSystem.SaveFileArgsForCall(0)
      Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
      Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

    })

  })

})
