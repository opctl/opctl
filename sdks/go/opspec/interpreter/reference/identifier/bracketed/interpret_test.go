package bracketed

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/bracketed/item"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("ref doesn't start w/ '['", func() {
		It("should return expected result", func() {
			/* arrange */
			providedRef := "dummyRef"

			expectedErr := fmt.Errorf("unable to interpret '%v': expected '['", providedRef)

			/* act */
			_, _, actualErr := Interpret(
				providedRef,
				nil,
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Context("ref doesn't contain ']", func() {
		It("should return expected result", func() {
			/* arrange */
			providedRef := "[dummyRef"

			expectedErr := fmt.Errorf("unable to interpret '%v': expected ']'", providedRef)

			/* act */
			_, _, actualErr := Interpret(
				providedRef,
				nil,
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Context("coerceToArrayOrObjecter.CoerceToArrayOrObject errs", func() {

		It("should return expected results", func() {

			/* arrange */
			providedRef := "[]"
			providedData := model.Value{String: new(string)}

			/* act */
			_, _, actualErr := Interpret(
				providedRef,
				&providedData,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret '[]': unable to coerce string to object: unexpected end of JSON input"))
		})
	})
	Context("data is array", func() {
		Context("item.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */

				providedRefIdentifier := "dummyIdentifier"
				providedRef := fmt.Sprintf("[%v]", providedRefIdentifier)

				arrayValue := []interface{}{nil}
				providedData := model.Value{Array: &arrayValue}

				/* act */
				_, _, actualErr := Interpret(
					providedRef,
					&providedData,
				)

				/* assert */
				Expect(actualErr).To(MatchError("unable to interpret item: strconv.ParseInt: parsing \"dummyIdentifier\": invalid syntax"))
			})
		})
		Context("item.Interpret doesn't err", func() {

			It("should return expected result", func() {
				/* arrange */

				providedRefIdentifier := "0"
				providedRef := fmt.Sprintf("[%v]", providedRefIdentifier)

				arrayValue := []interface{}{"item"}
				providedData := model.Value{Array: &arrayValue}

				expectedValue, err := item.Interpret(providedRefIdentifier, providedData)
				if err != nil {
					panic(err)
				}

				/* act */
				actualRefRemainder, actualData, actualErr := Interpret(
					providedRef,
					&providedData,
				)

				/* assert */
				Expect(actualRefRemainder).To(BeEmpty())
				Expect(*actualData).To(Equal(*expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("data is Object", func() {
		Context("value.Construct errs", func() {
			It("should return expected result", func() {
				/* arrange */

				providedRefIdentifier := "dummyIdentifier"
				providedRef := fmt.Sprintf("[%v]", providedRefIdentifier)

				object := &map[string]interface{}{providedRefIdentifier: nil}
				providedData := model.Value{Object: object}

				/* act */
				_, _, actualErr := Interpret(
					providedRef,
					&providedData,
				)

				/* assert */
				Expect(actualErr).To(MatchError("unable to interpret property: unable to construct value: '<nil>' unexpected type"))
			})
		})
		Context("value.Construct doesn't err", func() {

			It("should return expected result", func() {
				/* arrange */

				providedRefIdentifier := "dummyIdentifier"
				providedRef := fmt.Sprintf("[%v]", providedRefIdentifier)

				object := &map[string]interface{}{providedRefIdentifier: "string"}
				providedData := model.Value{Object: object}

				expectedData, err := value.Construct((*providedData.Object)[providedRefIdentifier])
				if err != nil {
					panic(err)
				}

				/* act */
				actualRefRemainder, actualData, actualErr := Interpret(
					providedRef,
					&providedData,
				)

				/* assert */
				Expect(actualRefRemainder).To(BeEmpty())
				Expect(*actualData).To(Equal(*expectedData))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
