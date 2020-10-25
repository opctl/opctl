package predicates

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	predicateFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/predicate/fakes"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("should call predicateInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedPredicateSpecs := []*model.PredicateSpec{
				&model.PredicateSpec{Eq: new([]interface{})},
				&model.PredicateSpec{Ne: new([]interface{})},
			}

			providedScope := map[string]*model.Value{}

			fakePredicateInterpreter := new(predicateFakes.FakeInterpreter)
			fakePredicateInterpreter.InterpretReturns(
				true,
				nil,
			)

			objectUnderTest := _interpreter{
				predicateInterpreter: fakePredicateInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedPredicateSpecs,
				providedScope,
			)

			/* assert */
			actualPredicateSpecs0,
				actualScope0 := fakePredicateInterpreter.InterpretArgsForCall(0)

			Expect(actualPredicateSpecs0).To(Equal(providedPredicateSpecs[0]))
			Expect(actualScope0).To(Equal(providedScope))

			actualPredicateSpecs1,
				actualScope1 := fakePredicateInterpreter.InterpretArgsForCall(1)

			Expect(actualPredicateSpecs1).To(Equal(providedPredicateSpecs[1]))
			Expect(actualScope1).To(Equal(providedScope))
		})
		Context("predicateInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakePredicateInterpreter := new(predicateFakes.FakeInterpreter)

				expectedError := errors.New("expectedError")
				fakePredicateInterpreter.InterpretReturns(
					false,
					expectedError,
				)

				objectUnderTest := _interpreter{
					predicateInterpreter: fakePredicateInterpreter,
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					[]*model.PredicateSpec{
						&model.PredicateSpec{Eq: new([]interface{})},
					},
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("predicateInterpreter.Interpret returns true", func() {
			It("should return expected result", func() {
				/* arrange */
				fakePredicateInterpreter := new(predicateFakes.FakeInterpreter)

				fakePredicateInterpreter.InterpretReturns(
					true,
					nil,
				)

				objectUnderTest := _interpreter{
					predicateInterpreter: fakePredicateInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					[]*model.PredicateSpec{
						&model.PredicateSpec{Eq: new([]interface{})},
						&model.PredicateSpec{Ne: new([]interface{})},
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
				fakePredicateInterpreter := new(predicateFakes.FakeInterpreter)

				fakePredicateInterpreter.InterpretReturns(
					false,
					nil,
				)

				objectUnderTest := _interpreter{
					predicateInterpreter: fakePredicateInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					[]*model.PredicateSpec{
						&model.PredicateSpec{Eq: new([]interface{})},
						&model.PredicateSpec{Ne: new([]interface{})},
					},
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(BeFalse())
			})
		})
	})
})
