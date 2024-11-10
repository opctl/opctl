package object

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
				map[string]*ipld.Node{},
				"$()",
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret $() to object: unable to interpret '' as reference: '' not in scope"))

		})
	})
	Context("value.Interpret doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			objectData := map[string]interface{}{}
			/* act */
			actualObject, actualErr := Interpret(
				map[string]*ipld.Node{},
				"{}",
			)

			/* assert */

			Expect(*actualObject).To(Equal(ipld.Node{Object: &objectData}))
			Expect(actualErr).To(BeNil())
		})
	})
})
