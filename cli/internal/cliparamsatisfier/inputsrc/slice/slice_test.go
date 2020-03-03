package slice

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sliceInputSrc", func() {
	Context("ReadString()", func() {
		Context("args doesn't contain delimited entry w/ provided inputName", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := New([]string{}, "")

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString("nonExistentInputName")

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualOk).To(BeFalse())
			})
		})
		Context("args contains delimited entry w/ provided inputName", func() {
			It("should return value", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				sep := "="
				expectedValue := "dummyValue"

				objectUnderTest := New(
					[]string{
						fmt.Sprintf("%v%v%v", providedInputName, sep, expectedValue),
					},
					sep,
				)

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString(providedInputName)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualOk).To(BeTrue())
			})
			It("should return value only once", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				sep := "="
				expectedValue := "dummyValue"

				objectUnderTest := New(
					[]string{
						fmt.Sprintf("%v%v%v", providedInputName, sep, expectedValue),
					},
					sep,
				)

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
