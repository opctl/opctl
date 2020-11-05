package object

import (
	"errors"

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
			Expect(actualErr).To(Equal(errors.New("unable to interpret $() to object; error was unable to interpret '' as reference; '' not in scope")))

		})
	})
	Context("value.Interpret doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			objectData := map[string]interface{}{}
			/* act */
			actualObject, actualErr := Interpret(
				map[string]*model.Value{},
				"{}",
			)

			/* assert */

			Expect(*actualObject).To(Equal(model.Value{Object: &objectData}))
			Expect(actualErr).To(BeNil())
		})
	})
})
