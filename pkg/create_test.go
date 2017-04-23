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

var _ = Describe("pkg", func() {

	Context("Create", func() {

		It("should call FileSystem.MkdirAll with expected args", func() {
			/* arrange */
			providedPath := "dummyPath"
			expectedPerm := os.FileMode(0777)

			fakeFileSystem := new(fs.Fake)

			objectUnderTest := &pkg{
				fs:     fakeFileSystem,
				ioUtil: new(vioutil.Fake),
			}

			/* act */
			objectUnderTest.Create(providedPath, "", "")

			/* assert */
			actualPath, actualPerm := fakeFileSystem.MkdirAllArgsForCall(0)
			Expect(actualPath).To(Equal(providedPath))
			Expect(actualPerm).To(Equal(expectedPerm))
		})

		Context("when FileSystem.MkdirAll returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("AddDirError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.MkdirAllReturns(expectedError)

				objectUnderTest := &pkg{
					fs: fakeFileSystem,
				}

				/* act */
				actualError := objectUnderTest.Create("", "", "")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})
	})

	It("should call ioutil.WriteFile with expected args", func() {

		/* arrange */
		providedPath := "dummyPath"
		providedPkgName := "dummyPkgName"
		providedPkgDescription := "dummyPkgDescription"

		expectedPkgManifestBytes, err := yaml.Marshal(&model.PkgManifest{
			Description: providedPkgDescription,
			Name:        providedPkgName,
		})
		if nil != err {
			panic(err)
		}

		expectedPath := path.Join(providedPath, OpDotYmlFileName)
		expectedData := expectedPkgManifestBytes
		expectedPerms := os.FileMode(0777)

		fakeIOUtil := new(vioutil.Fake)

		objectUnderTest := &pkg{
			fs:     new(fs.Fake),
			ioUtil: fakeIOUtil,
		}

		/* act */
		objectUnderTest.Create(providedPath, providedPkgName, providedPkgDescription)

		/* assert */
		actualPath, actualData, actualPerms := fakeIOUtil.WriteFileArgsForCall(0)
		Expect(actualPath).To(Equal(expectedPath))
		Expect(actualData).To(Equal(expectedData))
		Expect(actualPerms).To(Equal(expectedPerms))
	})

})
