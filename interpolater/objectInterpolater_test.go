package interpolater

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Interpolate", func() {
	Context("passed input containing no var placeholders", func() {
		It("should return input unmodified", func() {
			/* arrange */
			providedVarName := "dummyVarName"
			providedVarValue := map[string]interface{}{"dummyProp1Name": "dummyProp1Value"}
			providedInput := "dummyInput"
			expectedResult := providedInput
			objectUnderTest := newObjectInterpolater()

			/* act */
			actualResult := objectUnderTest.Interpolate(providedInput, providedVarName, providedVarValue)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
	})
	Context("input containing placeholders not referencing varName", func() {
		It("should return input unmodified", func() {
			/* arrange */
			providedVarName := "dummyVarName"
			providedVarValue := map[string]interface{}{"dummyProp1Name": "dummyProp1Value"}
			providedInput := "dummyInput $(var1) $(var2)"
			expectedResult := providedInput

			objectUnderTest := newObjectInterpolater()

			/* act */
			actualResult := objectUnderTest.Interpolate(providedInput, providedVarName, providedVarValue)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
	})
	Context("input containing placeholders referencing varName", func() {
		It("should replace all placeholders referencing varName", func() {
			/* arrange */
			providedVarName := "dummyVarName"
			providedVarValue := map[string]interface{}{"dummyProp1Name": "dummyProp1Value"}
			providedVarValueBytes, err := json.Marshal(providedVarValue)
			if nil != err {
				Fail(err.Error())
			}

			providedInput := fmt.Sprintf(
				"dummyInput $(%v) $(%v)",
				providedVarName,
				providedVarName,
			)

			expectedResult := fmt.Sprintf(
				"dummyInput %v %v",
				string(providedVarValueBytes),
				string(providedVarValueBytes),
			)

			objectUnderTest := newObjectInterpolater()

			/* act */
			actualResult := objectUnderTest.Interpolate(providedInput, providedVarName, providedVarValue)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
		Context("and placeholders not referencing varName", func() {
			It("should replace only placeholders referencing varName", func() {
				/* arrange */
				providedVarName := "dummyVarName"
				providedVarValue := map[string]interface{}{"dummyProp1Name": "dummyProp1Value"}

				providedVarValueBytes, err := json.Marshal(providedVarValue)
				if nil != err {
					Fail(err.Error())
				}

				providedInput := fmt.Sprintf(
					"dummyInput $(%v) $(not%v) $(%v)",
					providedVarName,
					providedVarName,
					providedVarName,
				)

				expectedResult := fmt.Sprintf(
					"dummyInput %v $(not%v) %v",
					string(providedVarValueBytes),
					providedVarName,
					string(providedVarValueBytes),
				)

				objectUnderTest := newObjectInterpolater()

				/* act */
				actualResult := objectUnderTest.Interpolate(providedInput, providedVarName, providedVarValue)

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
	})
})
