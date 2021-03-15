package predicate

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("Eq Predicate", func() {
		It("should return expected result", func() {
			/* arrange */
			eqPredicate := []interface{}{
				"same",
				"same",
			}

			/* act */
			actualResult, actualError := Interpret(
				&model.PredicateSpec{
					Eq: &eqPredicate,
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(Equal(true))
			Expect(actualError).To(BeNil())
		})
	})
	Context("Exists Predicate", func() {
		It("should return expected result", func() {
			/* arrange */
			identifier := "identifier"
			existsPredicate := fmt.Sprintf("$(%s)", identifier)

			/* act */
			actualResult, actualError := Interpret(
				&model.PredicateSpec{
					Exists: &existsPredicate,
				},
				map[string]*model.Value{
					identifier: &model.Value{String: new(string)},
				},
			)

			/* assert */
			Expect(actualResult).To(Equal(true))
			Expect(actualError).To(BeNil())
		})
	})
	Context("Ne predicate", func() {
		It("should return expected result", func() {
			/* arrange */
			nePredicate := []interface{}{
				"not",
				"same",
			}

			/* act */
			actualResult, actualError := Interpret(
				&model.PredicateSpec{
					Ne: &nePredicate,
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(Equal(true))
			Expect(actualError).To(BeNil())
		})
	})
	Context("NotExists Predicate", func() {
		It("should return expected result", func() {
			/* arrange */
			identifier := "identifier"
			notExistsPredicate := fmt.Sprintf("$(%s)", identifier)

			/* act */
			actualResult, actualError := Interpret(
				&model.PredicateSpec{
					NotExists: &notExistsPredicate,
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(Equal(true))
			Expect(actualError).To(BeNil())
		})
	})
	Context("Unexpected predicate", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScgPredicate := &model.PredicateSpec{}

			expectedError := fmt.Errorf("unable to interpret predicate: predicate was unexpected type %+v", providedScgPredicate)

			/* act */
			_, actualError := Interpret(
				providedScgPredicate,
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualError).To(MatchError(expectedError))
		})
	})
})
