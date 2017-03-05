package pkg

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

var _ = Describe("_createCollection", func() {

	Context("Execute", func() {

		It("should call FileSystem.AddDir with expected args", func() {

			/* arrange */

			providedCreateCollectionReq := model.CreateCollectionReq{Path: "/dummy/path"}

			fakeFileSystem := new(fs.Fake)

			objectUnderTest := &pkg{
				fileSystem: fakeFileSystem,
				yaml:       format.NewYamlFormat(),
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

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.AddDirReturns(expectedError)

				objectUnderTest := &pkg{
					fileSystem: fakeFileSystem,
					yaml:       format.NewYamlFormat(),
				}

				/* act */
				actualError := objectUnderTest.CreateCollection(
					model.CreateCollectionReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		Context("when YamlFormat.From returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("ToError")

				fakeYamlFormat := new(format.Fake)
				fakeYamlFormat.FromReturns(nil, expectedError)

				objectUnderTest := &pkg{
					fileSystem: new(fs.Fake),
					yaml:       fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.CreateCollection(
					model.CreateCollectionReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		It("should call YamlFormat.From with expected collectionManifest", func() {

			/* arrange */
			expectedCollectionManifest := model.CollectionManifest{
				Manifest: model.Manifest{
					Description: "DummyDescription",
					Name:        "DummyName",
				},
			}

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
				reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedCollectionManifest))
				return
			}

			objectUnderTest := &pkg{
				fileSystem: new(fs.Fake),
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.CreateCollection(
				model.CreateCollectionReq{
					Description: expectedCollectionManifest.Description,
					Name:        expectedCollectionManifest.Name,
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

			fakeFileSystem := new(fs.Fake)

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

			objectUnderTest := &pkg{
				fileSystem: fakeFileSystem,
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.CreateCollection(
				model.CreateCollectionReq{Path: providedPath},
			)

			/* assert */
			actualSaveFilePathArg, actualSaveFileBytesArg := fakeFileSystem.SaveFileArgsForCall(0)
			Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
			Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

		})

	})

})
