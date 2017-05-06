package cliparamsatisfier

import (
	"github.com/golang-interfaces/vos"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("envVarInputSrc", func() {
	Context("NewEnvVarInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewEnvVarInputSrc()).To(Not(BeNil()))
		})
	})
	Context("Read()", func() {
		It("should call os.Getenv w/ expected args", func() {
			/* arrange */
			providedInputName := "dummyInputName"

			fakeOS := new(vos.Fake)
			objectUnderTest := envVarInputSrc{
				os:          fakeOS,
				readHistory: map[string]struct{}{},
			}

			/* act */
			objectUnderTest.Read(providedInputName)

			/* assert */
			Expect(fakeOS.GetenvArgsForCall(0)).To(Equal(providedInputName))
		})
		Context("os.Getenv returns empty value", func() {
			It("should return nil", func() {
				/* arrange */
				fakeOS := new(vos.Fake)

				objectUnderTest := envVarInputSrc{
					os:          fakeOS,
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue := objectUnderTest.Read("")

				/* assert */
				Expect(actualValue).To(BeNil())
			})
		})
		Context("os.Getenv returns non-empty value", func() {
			It("should return value", func() {
				/* arrange */
				expectedValue := "dummyValue"

				fakeOS := new(vos.Fake)
				fakeOS.GetenvReturns(expectedValue)

				objectUnderTest := envVarInputSrc{
					os:          fakeOS,
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue := objectUnderTest.Read("")

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
			})
			It("should return value only once", func() {
				/* arrange */
				expectedValue := "dummyValue"

				fakeOS := new(vos.Fake)
				fakeOS.GetenvReturns(expectedValue)

				objectUnderTest := envVarInputSrc{
					os:          fakeOS,
					readHistory: map[string]struct{}{},
				}

				/* act */
				actualValue1 := objectUnderTest.Read("")
				actualValue2 := objectUnderTest.Read("")

				/* assert */
				Expect(*actualValue1).To(Equal(expectedValue))
				Expect(actualValue2).To(BeNil())
			})
		})
	})
})
