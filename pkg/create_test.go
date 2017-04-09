package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/virtual-go/fs"
	"github.com/virtual-go/vioutil"
	"gopkg.in/yaml.v2"
	"os"
	"path"
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
				}

				/* act */
				actualError := objectUnderTest.Create(
					CreateReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})
	})

	It("should call ioutil.WriteFile with expected args", func() {

		/* arrange */
		providedCreateReq := CreateReq{
			Path:        "dummyPath",
			Description: "dummyDescription",
			Name:        "dummyName",
		}

		expectedPkgManifestBytes, err := yaml.Marshal(&model.PkgManifest{
			Description: providedCreateReq.Description,
			Name:        providedCreateReq.Name,
		})
		if nil != err {
			panic(err)
		}

		expectedPath := path.Join(providedCreateReq.Path, ManifestFileName)
		expectedData := expectedPkgManifestBytes
		expectedPerms := os.FileMode(0777)

		fakeIOUtil := new(vioutil.Fake)

		objectUnderTest := &pkg{
			fileSystem: new(fs.Fake),
			ioUtil:     fakeIOUtil,
		}

		/* act */
		objectUnderTest.Create(providedCreateReq)

		/* assert */
		actualPath, actualData, actualPerms := fakeIOUtil.WriteFileArgsForCall(0)
		Expect(actualPath).To(Equal(expectedPath))
		Expect(actualData).To(Equal(expectedData))
		Expect(actualPerms).To(Equal(expectedPerms))
	})

})
