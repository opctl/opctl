package cliparamsatisfier

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sliceInputSrc", func() {
	Context("Read()", func() {
		Context("args doesn't contain delimited entry w/ provided inputName", func() {
			It("should return nil", func() {
				/* arrange */
				objectUnderTest := NewSliceInputSrc([]string{}, "")

				/* act */
				actualValue := objectUnderTest.Read("nonExistentInputName")

				/* assert */
				Expect(actualValue).To(BeNil())
			})
		})
		Context("args contains delimited entry w/ provided inputName", func() {
			It("should return value", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				sep := "="
				expectedValue := "dummyValue"

				objectUnderTest := NewSliceInputSrc(
					[]string{
						fmt.Sprintf("%v%v%v", providedInputName, sep, expectedValue),
					},
					sep,
				)

				/* act */
				actualValue := objectUnderTest.Read(providedInputName)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
			})
			It("should return value only once", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				sep := "="
				expectedValue := "dummyValue"

				objectUnderTest := NewSliceInputSrc(
					[]string{
						fmt.Sprintf("%v%v%v", providedInputName, sep, expectedValue),
					},
					sep,
				)

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
