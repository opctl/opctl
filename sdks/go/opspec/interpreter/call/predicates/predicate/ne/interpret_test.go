package ne

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("str.Interpret errs", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			_, actualError := Interpret(
				[]interface{}{"$()"},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualError).To(MatchError("unable to interpret $() to string: unable to interpret '' as reference: '' not in scope"))
		})
	})
	Context("str.Interpret returns equal items", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualResult, _ := Interpret(
				[]interface{}{
					"expression0",
					"expression0",
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(BeFalse())
		})
	})
	Context("str.Interpret returns unequal items", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualResult, _ := Interpret(
				[]interface{}{
					"expression0",
					"expression1",
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(BeTrue())
		})
	})
})
