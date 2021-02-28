package ymlfile

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("ymlFileInputSrc", func() {
	Context("ReadString()", func() {
		wd, err := os.Getwd()
		if nil != err {
			panic(err)
		}
		argsYmlTestDataPath := filepath.Join(wd, "testdata/args.yml")
		Context("yml doesn't contain entry w/ provided inputName", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest, err := New(argsYmlTestDataPath)
				if err != nil {
					Fail(err.Error())
				}

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString("nonExistentInputName")

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualOk).To(BeFalse())
			})
		})
		Context("yml contains entry w/ provided inputName", func() {
			It("should return expected result", func() {
				/* arrange */
				providedInputName := "hello"
				expectedValue := "world"

				objectUnderTest, err := New(argsYmlTestDataPath)
				if err != nil {
					Fail(err.Error())
				}

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString(providedInputName)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualOk).To(BeTrue())
			})
			It("should return value only once", func() {
				/* arrange */
				providedInputName := "hello"
				expectedValue := "world"

				objectUnderTest, err := New(argsYmlTestDataPath)
				if err != nil {
					Fail(err.Error())
				}

				/* act */
				actualValue1, actualOk1 := objectUnderTest.ReadString(providedInputName)
				actualValue2, actualOk2 := objectUnderTest.ReadString(providedInputName)

				/* assert */
				Expect(*actualValue1).To(Equal(expectedValue))
				Expect(actualOk1).To(BeTrue())

				Expect(actualValue2).To(BeNil())
				Expect(actualOk2).To(BeFalse())
			})
		})
	})
})
