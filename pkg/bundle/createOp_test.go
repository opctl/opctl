package bundle

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
	"path"
	"reflect"
)

var _ = Describe("_createOp", func() {

	Context("Execute", func() {

		It("should call FileSystem.AddDir with expected args", func() {

			/* arrange */

			providedCreateOpReq := model.CreateOpReq{Path: "/dummy/path"}

			fakeFileSystem := new(fs.FakeFileSystem)

			objectUnderTest := &_bundle{
				fileSystem: fakeFileSystem,
				yaml:       format.NewYamlFormat(),
			}

			/* act */
			objectUnderTest.CreateOp(
				providedCreateOpReq,
			)

			/* assert */
			Expect(fakeFileSystem.AddDirArgsForCall(0)).To(Equal(providedCreateOpReq.Path))

		})

		Context("when FileSystem.AddDir returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("AddDirError")

				fakeFileSystem := new(fs.FakeFileSystem)
				fakeFileSystem.AddDirReturns(expectedError)

				objectUnderTest := &_bundle{
					fileSystem: fakeFileSystem,
					yaml:       format.NewYamlFormat(),
				}

				/* act */
				actualError := objectUnderTest.CreateOp(
					model.CreateOpReq{},
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
				fakeYamlFormat.FromReturns(nil, expectedError)

				objectUnderTest := &_bundle{
					fileSystem: new(fs.FakeFileSystem),
					yaml:       fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.CreateOp(
					model.CreateOpReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		It("should call YamlFormat.From with expected opManifest", func() {

			/* arrange */
			expectedOpManifest := model.OpManifest{
				Manifest: model.Manifest{
					Description: "DummyDescription",
					Name:        "DummyName",
				},
			}

			fakeYamlFormat := new(format.FakeFormat)
			fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
				reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedOpManifest))
				return
			}

			objectUnderTest := &_bundle{
				fileSystem: new(fs.FakeFileSystem),
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.CreateOp(
				model.CreateOpReq{
					Description: expectedOpManifest.Description,
					Name:        expectedOpManifest.Name,
				},
			)

			/* assert */
			actualOpManifest := fakeYamlFormat.FromArgsForCall(0)
			Expect(actualOpManifest).To(Equal(&expectedOpManifest))

		})

		It("should call FileSystem.SaveFile with expected args", func() {

			/* arrange */
			providedPath := "/dummy/op/path"
			expectedSaveFilePathArg := path.Join(providedPath, NameOfOpManifestFile)
			expectedSaveFileBytesArg := []byte{2, 3, 4}

			fakeFileSystem := new(fs.FakeFileSystem)

			fakeYamlFormat := new(format.FakeFormat)
			fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

			objectUnderTest := &_bundle{
				fileSystem: fakeFileSystem,
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.CreateOp(
				model.CreateOpReq{Path: providedPath},
			)

			/* assert */
			actualSaveFilePathArg, actualSaveFileBytesArg := fakeFileSystem.SaveFileArgsForCall(0)
			Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
			Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

		})

	})

})
