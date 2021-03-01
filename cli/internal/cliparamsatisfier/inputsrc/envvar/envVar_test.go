package envvar

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("envVarInputSrc", func() {
	Context("ReadString()", func() {
		Context("os.Getenv returns empty value", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := envVarInputSrc{
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString("DOESNT_EXIST")

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualOk).To(BeFalse())
			})
		})
		Context("os.Getenv returns non-empty value", func() {
			It("should return value", func() {
				/* arrange */
				envVarName := "DOES_EXIST"
				expectedValue := "dummyValue"

				os.Setenv(envVarName, expectedValue)

				objectUnderTest := envVarInputSrc{
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString(envVarName)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualOk).To(BeTrue())
			})
			It("should return value only once", func() {
				/* arrange */
				envVarName := "DOES_EXIST"
				expectedValue := "dummyValue"

				os.Setenv(envVarName, expectedValue)

				objectUnderTest := envVarInputSrc{
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue1, actualOk1 := objectUnderTest.ReadString(envVarName)
				actualValue2, actualOk2 := objectUnderTest.ReadString(envVarName)

				/* assert */
				Expect(*actualValue1).To(Equal(expectedValue))
				Expect(actualOk1).To(BeTrue())

				Expect(actualValue2).To(BeNil())
				Expect(actualOk2).To(BeFalse())
			})
		})
	})
})
