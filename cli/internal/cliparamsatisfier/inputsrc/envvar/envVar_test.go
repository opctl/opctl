package envvar

import (
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("envVarInputSrc", func() {
	Context("ReadString()", func() {
		It("should call os.Getenv w/ expected args", func() {
			/* arrange */
			providedInputName := "dummyInputName"

			fakeOS := new(ios.Fake)
			objectUnderTest := envVarInputSrc{
				os:          fakeOS,
				readHistory: map[string]struct{}{},
			}

			/* act */
			objectUnderTest.ReadString(providedInputName)

			/* assert */
			Expect(fakeOS.GetenvArgsForCall(0)).To(Equal(providedInputName))
		})
		Context("os.Getenv returns empty value", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeOS := new(ios.Fake)

				objectUnderTest := envVarInputSrc{
					os:          fakeOS,
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString("")

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualOk).To(BeFalse())
			})
		})
		Context("os.Getenv returns non-empty value", func() {
			It("should return value", func() {
				/* arrange */
				expectedValue := "dummyValue"

				fakeOS := new(ios.Fake)
				fakeOS.GetenvReturns(expectedValue)

				objectUnderTest := envVarInputSrc{
					os:          fakeOS,
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString("")

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualOk).To(BeTrue())
			})
			It("should return value only once", func() {
				/* arrange */
				expectedValue := "dummyValue"

				fakeOS := new(ios.Fake)
				fakeOS.GetenvReturns(expectedValue)

				objectUnderTest := envVarInputSrc{
					os:          fakeOS,
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue1, actualOk1 := objectUnderTest.ReadString("")
				actualValue2, actualOk2 := objectUnderTest.ReadString("")

				/* assert */
				Expect(*actualValue1).To(Equal(expectedValue))
				Expect(actualOk1).To(BeTrue())

				Expect(actualValue2).To(BeNil())
				Expect(actualOk2).To(BeFalse())
			})
		})
	})
})
