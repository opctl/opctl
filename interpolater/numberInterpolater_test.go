package interpolater

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interpolate", func() {
	Describe("passed input containing no var placeholders", func() {
		It("should return input unmodified", func() {
			/* arrange */
			providedVarName := "dummyVarName"
			providedVarValue := 1.0
			providedInput := "dummyInput"
			expectedResult := providedInput
			objectUnderTest := newNumberInterpolater()

			/* act */
			actualResult := objectUnderTest.Interpolate(providedInput, providedVarName, providedVarValue)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
	})
	Describe("input containing placeholders not referencing varName", func() {
		It("should return input unmodified", func() {
			/* arrange */
			providedVarName := "dummyVarName"
			providedVarValue := 1.0
			providedInput := "dummyInput $(var1) $(var2)"
			expectedResult := providedInput

			objectUnderTest := newNumberInterpolater()

			/* act */
			actualResult := objectUnderTest.Interpolate(providedInput, providedVarName, providedVarValue)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
	})
	Describe("input containing placeholders referencing varName", func() {
		It("should replace all placeholders referencing varName", func() {
			/* arrange */
			providedVarName := "dummyVarName"
			providedVarValue := 1.0
			providedInput := fmt.Sprintf("dummyInput $(%v) $(%v)", providedVarName, providedVarName)
			expectedResult := fmt.Sprintf("dummyInput %v %v", providedVarValue, providedVarValue)

			objectUnderTest := newNumberInterpolater()

			/* act */
			actualResult := objectUnderTest.Interpolate(providedInput, providedVarName, providedVarValue)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
		Describe("and placeholders not referencing varName", func() {
			It("should replace only placeholders referencing varName", func() {
				/* arrange */
				providedVarName := "dummyVarName"
				providedVarValue := 1.0
				providedInput := fmt.Sprintf("dummyInput $(%v) $(not%v) $(%v)", providedVarName, providedVarName, providedVarName)
				expectedResult := fmt.Sprintf("dummyInput %v $(not%v) %v", providedVarValue, providedVarName, providedVarValue)

				objectUnderTest := newNumberInterpolater()

				/* act */
				actualResult := objectUnderTest.Interpolate(providedInput, providedVarName, providedVarValue)

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
	})
})
