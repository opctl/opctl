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

var _ = Describe("_createCollection", func() {

  Context("Execute", func() {

    It("should call FileSystem.AddDir with expected args", func() {

      /* arrange */

      providedCreateCollectionReq := models.CreateCollectionReq{Path:"/dummy/path"}

      fakeFileSystem := new(fs.FakeFileSystem)

      objectUnderTest := &_bundle{
        fileSystem: fakeFileSystem,
        yaml: format.NewYamlFormat(),
      }

      /* act */
      objectUnderTest.CreateCollection(
        providedCreateCollectionReq,
      )

      /* assert */
      Expect(fakeFileSystem.AddDirArgsForCall(0)).To(Equal(providedCreateCollectionReq.Path))

    })

    Context("when FileSystem.AddDir returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("AddDirError")

        fakeFileSystem := new(fs.FakeFileSystem)
        fakeFileSystem.AddDirReturns(expectedError)

        objectUnderTest := &_bundle{
          fileSystem: fakeFileSystem,
          yaml: format.NewYamlFormat(),
        }

        /* act */
        actualError := objectUnderTest.CreateCollection(
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

        fakeYamlFormat := new(format.FakeFormat)
        fakeYamlFormat.FromReturns(nil, expectedError)

        objectUnderTest := &_bundle{
          fileSystem: new(fs.FakeFileSystem),
          yaml: fakeYamlFormat,
        }

        /* act */
        actualError := objectUnderTest.CreateCollection(
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

      fakeYamlFormat := new(format.FakeFormat)
      fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedCollectionManifest))
        return
      }

      objectUnderTest := &_bundle{
        fileSystem: new(fs.FakeFileSystem),
        yaml: fakeYamlFormat,
      }

      /* act */
      objectUnderTest.CreateCollection(
        models.CreateCollectionReq{
          Description:expectedCollectionManifest.Description,
          Name:expectedCollectionManifest.Name,
        },
      )

      /* assert */
      actualCollectionManifest := fakeYamlFormat.FromArgsForCall(0)
      Expect(actualCollectionManifest).To(Equal(&expectedCollectionManifest))

    })

    It("should call FileSystem.SaveFile with expected args", func() {

      /* arrange */
      providedPath := "/dummy/op/path"
      expectedSaveFilePathArg := path.Join(providedPath, NameOfCollectionManifestFile)
      expectedSaveFileBytesArg := []byte{2, 3, 4}

      fakeFileSystem := new(fs.FakeFileSystem)

      fakeYamlFormat := new(format.FakeFormat)
      fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

      objectUnderTest := &_bundle{
        fileSystem: fakeFileSystem,
        yaml: fakeYamlFormat,
      }

      /* act */
      objectUnderTest.CreateCollection(
        models.CreateCollectionReq{Path:providedPath},
      )

      /* assert */
      actualSaveFilePathArg, actualSaveFileBytesArg := fakeFileSystem.SaveFileArgsForCall(0)
      Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
      Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

    })

  })

})
