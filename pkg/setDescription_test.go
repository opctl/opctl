package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/virtual-go/vioutil"
	"os"
	"path"
	"reflect"
)

var _ = Describe("_setDescription", func() {

	Context("Execute", func() {

		Context("when FileSystem.ReadFile returns an error", func() {

			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("GetBytesOfFileError")

				fakeIOUtil := new(vioutil.Fake)
				fakeIOUtil.ReadFileReturns(nil, expectedError)

				objectUnderTest := &pkg{
					ioUtil: fakeIOUtil,
					yaml:   format.NewYamlFormat(),
				}

				/* act */
				actualError := objectUnderTest.SetDescription(
					SetDescriptionReq{},
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

				objectUnderTest := &pkg{
					ioUtil: new(vioutil.Fake),
					yaml:   fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.SetDescription(
					SetDescriptionReq{},
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

				objectUnderTest := &pkg{
					ioUtil: new(vioutil.Fake),
					yaml:   fakeYamlFormat,
				}

				/* act */
				actualError := objectUnderTest.SetDescription(
					SetDescriptionReq{},
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

			objectUnderTest := &pkg{
				ioUtil: new(vioutil.Fake),
				yaml:   fakeYamlFormat,
			}

			/* act */
			objectUnderTest.SetDescription(
				SetDescriptionReq{Description: expectedPackageManifestView.Description},
			)

			/* assert */
			actualPackageManifestView := fakeYamlFormat.FromArgsForCall(0)
			Expect(actualPackageManifestView).To(Equal(&expectedPackageManifestView))

		})

		It("should call ioutil.WriteFile with expected args", func() {

			/* arrange */
			providedPath := "/dummy/op/path"
			expectedPath := path.Join(providedPath, NameOfPkgManifestFile)
			expectedBytes := []byte{2, 3, 4}
			expectedPerms := os.FileMode(0777)

			fakeIOUtil := new(vioutil.Fake)

			fakeYamlFormat := new(format.Fake)
			fakeYamlFormat.FromReturns(expectedBytes, nil)

			objectUnderTest := &pkg{
				ioUtil: fakeIOUtil,
				yaml:   fakeYamlFormat,
			}

			/* act */
			objectUnderTest.SetDescription(
				SetDescriptionReq{Path: providedPath},
			)

			/* assert */
			actualPath, actualBytes, actualPerms := fakeIOUtil.WriteFileArgsForCall(0)
			Expect(actualPath).To(Equal(expectedPath))
			Expect(actualBytes).To(Equal(expectedBytes))
			Expect(actualPerms).To(Equal(expectedPerms))

		})

	})

})
