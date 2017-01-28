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

var _ = Describe("_setOpDescription", func() {

	Context("Execute", func() {

		Context("when FileSystem.GetBytesOfFile returns an error", func() {

			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("GetBytesOfFileError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.GetBytesOfFileReturns(nil, expectedError)

				objectUnderTest := &_bundle{
					fileSystem: fakeFileSystem,
					yaml:       format.NewYamlFormat(),
				}

				/* act */
				actualError := objectUnderTest.SetOpDescription(
					model.SetOpDescriptionReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})

		})

		Context("when YamlFormat.From returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("FromError")

				fakeYamlFormat := new(format.Fake)
				fakeYamlFormat.ToReturns(expectedError)

				objectUnderTest := &_bundle{
					fileSystem: new(fs.Fake),
					yaml:       fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.SetOpDescription(
					model.SetOpDescriptionReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		Context("when YamlFormat.To returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("ToError")

				fakeYamlFormat := new(format.Fake)
				fakeYamlFormat.FromReturns(nil, expectedError)

				objectUnderTest := &_bundle{
					fileSystem: new(fs.Fake),
					yaml:       fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.SetOpDescription(
					model.SetOpDescriptionReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		It("should call YamlFormat.From with expected opManifest", func() {

			/* arrange */
			expectedOpManifest := model.OpManifest{
				Manifest: model.Manifest{
					Name:        "dummyName",
					Description: "dummyDescription",
					Version:     "dummyVersion",
				},
			}

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
				reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedOpManifest))
				return
			}

			objectUnderTest := &_bundle{
				fileSystem: new(fs.Fake),
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.SetOpDescription(
				model.SetOpDescriptionReq{Description: expectedOpManifest.Description},
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

			fakeFileSystem := new(fs.Fake)

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

			objectUnderTest := &_bundle{
				fileSystem: fakeFileSystem,
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.SetOpDescription(
				model.SetOpDescriptionReq{PathToOp: providedPathToOp},
			)

			/* assert */
			actualSaveFilePathArg, actualSaveFileBytesArg := fakeFileSystem.SaveFileArgsForCall(0)
			Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
			Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

		})

	})

})
