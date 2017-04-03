package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/virtual-go/fs"
	"github.com/virtual-go/vioutil"
	"os"
	"path"
	"reflect"
)

var _ = Describe("_create", func() {

	Context("Execute", func() {

		It("should call FileSystem.MkdirAll with expected args", func() {
			/* arrange */
			providedCreateReq := CreateReq{Path: "/dummy/path"}
			expectedPerm := os.FileMode(0777)

			fakeFileSystem := new(fs.Fake)

			objectUnderTest := &pkg{
				fileSystem: fakeFileSystem,
				ioUtil:     new(vioutil.Fake),
				yaml:       format.NewYamlFormat(),
			}

			/* act */
			objectUnderTest.Create(
				providedCreateReq,
			)

			/* assert */
			actualPath, actualPerm := fakeFileSystem.MkdirAllArgsForCall(0)
			Expect(actualPath).To(Equal(providedCreateReq.Path))
			Expect(actualPerm).To(Equal(expectedPerm))
		})

		Context("when FileSystem.MkdirAll returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("AddDirError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.MkdirAllReturns(expectedError)

				objectUnderTest := &pkg{
					fileSystem: fakeFileSystem,
					yaml:       format.NewYamlFormat(),
				}

				/* act */
				actualError := objectUnderTest.Create(
					CreateReq{},
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

				objectUnderTest := &pkg{
					fileSystem: new(fs.Fake),
					yaml:       fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.Create(
					CreateReq{},
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

			objectUnderTest := &pkg{
				fileSystem: new(fs.Fake),
				ioUtil:     new(vioutil.Fake),
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.Create(
				CreateReq{
					Description: expectedPackageManifestView.Description,
					Name:        expectedPackageManifestView.Name,
				},
			)

			/* assert */
			actualPackageManifestView := fakeYamlFormat.FromArgsForCall(0)
			Expect(actualPackageManifestView).To(Equal(&expectedPackageManifestView))

		})

		It("should call ioutil.WriteFile with expected args", func() {

			/* arrange */
			providedPath := "/dummy/op/path"
			expectedPath := path.Join(providedPath, NameOfPkgManifestFile)
			expectedData := []byte{2, 3, 4}
			expectedPerms := os.FileMode(0777)

			fakeIOUtil := new(vioutil.Fake)

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.FromReturns(expectedData, nil)

			objectUnderTest := &pkg{
				fileSystem: new(fs.Fake),
				ioUtil:     fakeIOUtil,
				yaml:       fakeYamlFormat,
			}

			/* act */
			objectUnderTest.Create(
				CreateReq{Path: providedPath},
			)

			/* assert */
			actualPath, actualData, actualPerms := fakeIOUtil.WriteFileArgsForCall(0)
			Expect(actualPath).To(Equal(expectedPath))
			Expect(actualData).To(Equal(expectedData))
			Expect(actualPerms).To(Equal(expectedPerms))
		})

	})

})
