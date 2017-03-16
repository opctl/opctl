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

var _ = Describe("_createPackage", func() {

	Context("Execute", func() {

		It("should call FileSystem.AddDir with expected args", func() {

			/* arrange */

			providedCreatePackageReq := model.CreatePackageReq{Path: "/dummy/path"}

			fakeFileSystem := new(fs.Fake)

			objectUnderTest := &managePackages{
				fileSystem: fakeFileSystem,
				yaml:       format.NewYamlFormat(),
			}

			/* act */
			objectUnderTest.CreatePackage(
				providedCreatePackageReq,
			)

			/* assert */
			Expect(fakeFileSystem.AddDirArgsForCall(0)).To(Equal(providedCreatePackageReq.Path))

		})

		Context("when FileSystem.AddDir returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("AddDirError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.AddDirReturns(expectedError)

				objectUnderTest := &managePackages{
					fileSystem: fakeFileSystem,
					yaml:       format.NewYamlFormat(),
				}

				/* act */
				actualError := objectUnderTest.CreatePackage(
					model.CreatePackageReq{},
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
				fakeYamlFormat.FromReturns(nil, expectedError)

				objectUnderTest := &managePackages{
					fileSystem: new(fs.Fake),
					yaml:       fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.CreatePackage(
					model.CreatePackageReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		It("should call YamlFormat.From with expected packageManifestView", func() {

			/* arrange */
			expectedPackageManifestView := model.PackageManifestView{
				Description: "DummyDescription",
				Name:        "DummyName",
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
			objectUnderTest.CreatePackage(
				model.CreatePackageReq{
					Description: expectedPackageManifestView.Description,
					Name:        expectedPackageManifestView.Name,
				},
			)

			/* assert */
			actualPackageManifestView := fakeYamlFormat.FromArgsForCall(0)
			Expect(actualPackageManifestView).To(Equal(&expectedPackageManifestView))

		})

		It("should call FileSystem.SaveFile with expected args", func() {

			/* arrange */
			providedPath := "/dummy/op/path"
			expectedSaveFilePathArg := path.Join(providedPath, NameOfPackageManifestFile)
			expectedSaveFileBytesArg := []byte{2, 3, 4}

			fakeFileSystem := new(fs.Fake)

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.FromReturns(expectedSaveFileBytesArg, nil)

			objectUnderTest := &managePackages{
				fileSystem: fakeFileSystem,
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.CreatePackage(
				model.CreatePackageReq{Path: providedPath},
			)

			/* assert */
			actualSaveFilePathArg, actualSaveFileBytesArg := fakeFileSystem.SaveFileArgsForCall(0)
			Expect(actualSaveFilePathArg).To(Equal(expectedSaveFilePathArg))
			Expect(actualSaveFileBytesArg).To(Equal(expectedSaveFileBytesArg))

		})

	})

})
