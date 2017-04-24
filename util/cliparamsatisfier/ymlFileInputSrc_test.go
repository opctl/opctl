package cliparamsatisfier

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/virtual-go/vioutil"
	"gopkg.in/yaml.v2"
)

var _ = Describe("ymlFileInputSrc", func() {
	Context("NewYMLFileInputSrc", func() {
		It("should call ioutil.ReadFile w/ expected args", func() {
			/* arrange */
			providedFilePath := "dummyFilePath"
			providedIOUtil := new(vioutil.Fake)
			providedIOUtil.ReadFileReturns(nil, errors.New(""))

			/* act */
			NewYMLFileInputSrc(providedFilePath, providedIOUtil)

			/* assert */
			Expect(providedIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFilePath))
		})
	})
	Context("Read()", func() {
		Context("yml doesn't contain entry w/ provided inputName", func() {
			It("should return nil", func() {
				/* arrange */
				argMap := map[string]string{}
				ymlBytes, err := yaml.Marshal(argMap)
				if nil != err {
					panic(err)
				}
				providedIOUtil := new(vioutil.Fake)
				providedIOUtil.ReadFileReturns(ymlBytes, nil)

				objectUnderTest := NewYMLFileInputSrc("", providedIOUtil)

				/* act */
				actualValue := objectUnderTest.Read("nonExistentInputName")

				/* assert */
				Expect(actualValue).To(BeNil())
			})
		})
		Context("yml contains entry w/ provided inputName", func() {
			It("should return value", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				expectedValue := "dummyValue"
				argMap := map[string]string{
					providedInputName: expectedValue,
				}
				ymlBytes, err := yaml.Marshal(argMap)
				if nil != err {
					panic(err)
				}
				providedIOUtil := new(vioutil.Fake)
				providedIOUtil.ReadFileReturns(ymlBytes, nil)

				objectUnderTest := NewYMLFileInputSrc("", providedIOUtil)

				/* act */
				actualValue := objectUnderTest.Read(providedInputName)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
			})
			It("should return value only once", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				expectedValue := "dummyValue"
				argMap := map[string]string{
					providedInputName: expectedValue,
				}
				ymlBytes, err := yaml.Marshal(argMap)
				if nil != err {
					panic(err)
				}
				providedIOUtil := new(vioutil.Fake)
				providedIOUtil.ReadFileReturns(ymlBytes, nil)

				objectUnderTest := NewYMLFileInputSrc("", providedIOUtil)

				/* act */
				actualValue1 := objectUnderTest.Read(providedInputName)
				actualValue2 := objectUnderTest.Read(providedInputName)

				/* assert */
				Expect(*actualValue1).To(Equal(expectedValue))
				Expect(actualValue2).To(BeNil())
			})
		})
	})
})
