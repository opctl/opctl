package loopable

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("array.Interpret doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			identifier := "identifier"
			arrayValue := []interface{}{"item"}

			providedScope := map[string]*model.Value{
				identifier: &model.Value{
					Array: &arrayValue,
				},
			}

			/* act */
			actualResult, actualErr := Interpret(
				fmt.Sprintf("$(%s)", identifier),
				providedScope,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(*actualResult).To(Equal(*providedScope[identifier]))
		})
	})
	Context("array.Interpret errs", func() {
		It("should return expected result", func() {
			/* arrange */
			identifier := "identifier"
			objectValue := map[string]interface{}{
				"key": "value",
			}

			providedScope := map[string]*model.Value{
				identifier: &model.Value{
					Object: &objectValue,
				},
			}

			/* act */
			actualResult, actualErr := Interpret(
				fmt.Sprintf("$(%s)", identifier),
				providedScope,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(*actualResult).To(Equal(*providedScope[identifier]))
		})
	})
})
