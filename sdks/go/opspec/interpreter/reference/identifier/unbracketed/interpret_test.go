package unbracketed

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"
)

var _ = Context("Interpret", func() {
	Context("coerce.ToObject errs", func() {
		It("should return expected result", func() {

			/* arrange */
			providedRef := "dummyRef"
			providedData := model.Value{String: new(string)}

			/* act */
			_, _, actualErr := Interpret(
				providedRef,
				&providedData,
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to interpret 'dummyRef'; error was unable to coerce string to object; error was unexpected end of JSON input")))
		})
	})
	Context("coerce.ToObject doesn't err", func() {

		Context("identifier not in object", func() {
			It("should return expected result", func() {

				/* arrange */
				identifier := "identifier"
				providedRef := fmt.Sprintf("%s.", identifier)

				objectData := map[string]interface{}{}

				expectedErr := fmt.Errorf("unable to interpret '%v'; '%v' doesn't exist", providedRef, identifier)

				/* act */
				_, _, actualErr := Interpret(
					providedRef,
					&model.Value{Object: &objectData},
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("identifier in object", func() {
			Context("value.Construct errs", func() {

				It("should return expected result", func() {

					/* arrange */
					identifier := "identifier"

					providedRef := fmt.Sprintf("%s.", identifier)

					objectData := map[string]interface{}{identifier: nil}

					/* act */
					_, _, actualErr := Interpret(
						providedRef,
						&model.Value{Object: &objectData},
					)

					/* assert */
					Expect(actualErr).To(Equal(errors.New("unable to interpret 'identifier.'; error was unable to construct value; '<nil>' unexpected type")))
				})
			})
			Context("value.Construct doesn't err", func() {

				It("should return expected result", func() {

					/* arrange */
					identifier := "identifier"

					providedRef := fmt.Sprintf("%s", identifier)

					objectData := map[string]interface{}{identifier: "dummyValue"}

					expectedValue, err := value.Construct(objectData[identifier])
					if nil != err {
						panic(err)
					}

					/* act */
					actualRefRemainder, actualValue, actualErr := Interpret(
						providedRef,
						&model.Value{Object: &objectData},
					)

					/* assert */
					Expect(actualRefRemainder).To(BeEmpty())
					Expect(*actualValue).To(Equal(*expectedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
