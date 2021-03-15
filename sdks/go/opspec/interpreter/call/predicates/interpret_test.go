package predicates

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("predicateInterpreter.Interpret errs", func() {
		It("should return expected result", func() {
			/* arrange */
			eqPredicate := []interface{}{
				"$()",
			}
			/* act */
			_, actualError := Interpret(
				[]*model.PredicateSpec{
					{Eq: &eqPredicate},
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualError).To(MatchError("unable to interpret $() to string: unable to interpret '' as reference: '' not in scope"))
		})
	})
	Context("predicateInterpreter.Interpret returns true", func() {
		It("should return expected result", func() {
			/* arrange */
			eqPredicate := []interface{}{
				true,
				true,
			}

			/* act */
			actualResult, _ := Interpret(
				[]*model.PredicateSpec{
					{Eq: &eqPredicate},
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(BeTrue())
		})
	})
	Context("predicateInterpreter.Interpret returns false", func() {
		It("should return expected result", func() {
			/* arrange */
			eqPredicate := []interface{}{
				true,
				false,
			}

			/* act */
			actualResult, _ := Interpret(
				[]*model.PredicateSpec{
					{Eq: &eqPredicate},
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(BeFalse())
		})
	})
})
