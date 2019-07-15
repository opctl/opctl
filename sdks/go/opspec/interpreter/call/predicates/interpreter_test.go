package predicates

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/predicates/predicate"
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
			providedSCGPredicates := []*model.SCGPredicate{
				&model.SCGPredicate{Eq: new([]interface{})},
				&model.SCGPredicate{Ne: new([]interface{})},
			}

			providedOpHandle := new(data.FakeHandle)
			providedScope := map[string]*model.Value{}

			fakePredicateInterpreter := new(predicate.FakeInterpreter)
			fakePredicateInterpreter.InterpretReturns(
				true,
				nil,
			)

			objectUnderTest := _interpreter{
				predicateInterpreter: fakePredicateInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedOpHandle,
				providedSCGPredicates,
				providedScope,
			)

			/* assert */
			actualOpHandle0,
				actualSCGPredicates0,
				actualScope0 := fakePredicateInterpreter.InterpretArgsForCall(0)

			Expect(actualOpHandle0).To(Equal(providedOpHandle))
			Expect(actualSCGPredicates0).To(Equal(providedSCGPredicates[0]))
			Expect(actualScope0).To(Equal(providedScope))

			actualOpHandle1,
				actualSCGPredicates1,
				actualScope1 := fakePredicateInterpreter.InterpretArgsForCall(1)

			Expect(actualOpHandle1).To(Equal(providedOpHandle))
			Expect(actualSCGPredicates1).To(Equal(providedSCGPredicates[1]))
			Expect(actualScope1).To(Equal(providedScope))
		})
		Context("predicateInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakePredicateInterpreter := new(predicate.FakeInterpreter)

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
					new(data.FakeHandle),
					[]*model.SCGPredicate{
						&model.SCGPredicate{Eq: new([]interface{})},
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
				fakePredicateInterpreter := new(predicate.FakeInterpreter)

				fakePredicateInterpreter.InterpretReturns(
					true,
					nil,
				)

				objectUnderTest := _interpreter{
					predicateInterpreter: fakePredicateInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					new(data.FakeHandle),
					[]*model.SCGPredicate{
						&model.SCGPredicate{Eq: new([]interface{})},
						&model.SCGPredicate{Ne: new([]interface{})},
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
				fakePredicateInterpreter := new(predicate.FakeInterpreter)

				fakePredicateInterpreter.InterpretReturns(
					false,
					nil,
				)

				objectUnderTest := _interpreter{
					predicateInterpreter: fakePredicateInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					new(data.FakeHandle),
					[]*model.SCGPredicate{
						&model.SCGPredicate{Eq: new([]interface{})},
						&model.SCGPredicate{Ne: new([]interface{})},
					},
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(BeFalse())
			})
		})
	})
})
