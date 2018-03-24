package op

import (
	"errors"
	"github.com/opspec-io/sdk-golang/op/dotyml"
	"os"
	"path/filepath"

	"github.com/chrisdostert/yaml"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("Creator", func() {
	Context("NewCreator", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewCreator()).Should(Not(BeNil()))
		})
	})
	Context("Create", func() {

		It("should call os.MkdirAll with expected args", func() {
			/* arrange */
			providedPath := "dummyPath"
			expectedPerm := os.FileMode(0777)

			fakeOS := new(ios.Fake)

			objectUnderTest := _creator{
				os:     fakeOS,
				ioUtil: new(iioutil.Fake),
			}

			/* act */
			objectUnderTest.Create(providedPath, "", "")

			/* assert */
			actualPath, actualPerm := fakeOS.MkdirAllArgsForCall(0)
			Expect(actualPath).To(Equal(providedPath))
			Expect(actualPerm).To(Equal(expectedPerm))
		})

		Context("when os.MkdirAll returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("AddDirError")

				fakeOS := new(ios.Fake)
				fakeOS.MkdirAllReturns(expectedError)

				objectUnderTest := _creator{
					os: fakeOS,
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

		expectedOpDotYmlBytes, err := yaml.Marshal(&model.OpDotYml{
			Description: providedPkgDescription,
			Name:        providedPkgName,
		})
		if nil != err {
			panic(err)
		}

		expectedPath := filepath.Join(providedPath, dotyml.FileName)
		expectedData := expectedOpDotYmlBytes
		expectedPerms := os.FileMode(0777)

		fakeIOUtil := new(iioutil.Fake)

		objectUnderTest := _creator{
			os:     new(ios.Fake),
			ioUtil: fakeIOUtil,
		}

		/* act */
		objectUnderTest.Create(
			providedPath,
			providedPkgName,
			providedPkgDescription,
		)

		/* assert */
		actualPath, actualData, actualPerms := fakeIOUtil.WriteFileArgsForCall(0)
		Expect(actualPath).To(Equal(expectedPath))
		Expect(actualData).To(Equal(expectedData))
		Expect(actualPerms).To(Equal(expectedPerms))
	})

})
