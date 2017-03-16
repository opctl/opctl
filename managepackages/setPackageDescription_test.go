package managepackages

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
	"path"
	"reflect"
)

var _ = Describe("_setPackageDescription", func() {

	Context("Execute", func() {

		Context("when FileSystem.GetBytesOfFile returns an error", func() {

			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("GetBytesOfFileError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.GetBytesOfFileReturns(nil, expectedError)

				objectUnderTest := &managePackages{
					fileSystem: fakeFileSystem,
					yaml:       format.NewYamlFormat(),
				}

				/* act */
				actualError := objectUnderTest.SetPackageDescription(
					model.SetPackageDescriptionReq{},
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

				objectUnderTest := &managePackages{
					fileSystem: new(fs.Fake),
					yaml:       fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.SetPackageDescription(
					model.SetPackageDescriptionReq{},
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

				objectUnderTest := &managePackages{
					fileSystem: new(fs.Fake),
					yaml:       fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.SetPackageDescription(
					model.SetPackageDescriptionReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		It("should call YamlFormat.From with expected packageManifestView", func() {

			/* arrange */
			expectedPackageManifestView := model.PackageManifestView{
				Name:        "dummyName",
				Description: "dummyDescription",
				Version:     "dummyVersion",
			}

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {
				reflect.ValueOf(out).Elem().Set(reflect.ValueOf(expectedPackageManifestView))
				return
			}

			objectUnderTest := &managePackages{
				fileSystem: new(fs.Fake),
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.SetPackageDescription(
				model.SetPackageDescriptionReq{Description: expectedPackageManifestView.Description},
			)

			/* assert */
			actualPackageManifestView := fakeYamlFormat.FromArgsForCall(0)
			Expect(actualPackageManifestView).To(Equal(&expectedPackageManifestView))

		})

		It("should call FileSystem.SaveFile with expected args", func() {

			/* arrange */
			providedPathToOp := "/dummy/op/path"
			expectedSaveFilePathArg := path.Join(providedPathToOp, NameOfPackageManifestFile)
			expectedSaveFileBytesArg := []byte{2, 3, 4}

			fakeFileSystem := new(fs.Fake)

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

			objectUnderTest := &managePackages{
				fileSystem: fakeFileSystem,
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.SetPackageDescription(
				model.SetPackageDescriptionReq{PathToOp: providedPathToOp},
			)

			/* assert */
			actualSaveFilePathArg, actualSaveFileBytesArg := fakeFileSystem.SaveFileArgsForCall(0)
			Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
			Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

		})

	})

})
