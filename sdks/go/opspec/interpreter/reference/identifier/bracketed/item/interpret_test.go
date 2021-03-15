package item

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"
)

var _ = Context("Interpret", func() {
	Context("parseIndexer.ParseIndex errs", func() {
		It("should return expected result", func() {
			/* arrange */

			/* act */
			_, actualErr := Interpret(
				"dummyIndexString",
				model.Value{Array: new([]interface{})},
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret item: strconv.ParseInt: parsing \"dummyIndexString\": invalid syntax"))
		})
	})
	Context("parseIndexer.ParseIndex doesn't err", func() {

		Context("value.Construct errs", func() {

			It("should return expected result", func() {
				/* arrange */

				arrayData := &[]interface{}{nil}
				providedData := model.Value{Array: arrayData}

				/* act */
				_, actualErr := Interpret(
					"0",
					providedData,
				)

				/* assert */
				Expect(actualErr).To(MatchError("unable to interpret item: unable to construct value: '<nil>' unexpected type"))
			})
		})
		Context("value.Construct doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */

				arrayData := &[]interface{}{"item"}
				providedData := model.Value{Array: arrayData}

				expectedValue, err := value.Construct((*providedData.Array)[0])
				if nil != err {
					panic(err)
				}

				/* act */
				actualItemValue, actualErr := Interpret(
					"0",
					providedData,
				)

				/* assert */
				Expect(*actualItemValue).To(Equal(*expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
