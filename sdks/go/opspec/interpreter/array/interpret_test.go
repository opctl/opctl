package array

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("value.Interpret errs", func() {
		It("should return expected err", func() {
			/* arrange */
			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{},
				"$()",
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret $() to array: unable to interpret '' as reference: '' not in scope"))

		})
	})
	Context("value.Interpret doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			arrayValue := []interface{}{}

			/* act */
			actualArray, actualErr := Interpret(
				map[string]*model.Value{},
				"[]",
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(*actualArray).To(Equal(model.Value{Array: &arrayValue}))
		})
	})
})
