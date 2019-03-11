package predicate

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/if/predicates/predicate/eq"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/if/predicates/predicate/ne"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("Eq Predicate", func() {
			It("should call eqInterpreter.Interpret w/ expected args & return result", func() {
				/* arrange */
				providedOpHandle := new(data.FakeHandle)

				providedScgPredicate := &model.SCGPredicate{
					Eq: []interface{}{},
				}

				providedScope := map[string]*model.Value{}

				fakeEqInterpreter := new(eq.FakeInterpreter)
				expectedResult := true
				expectedError := errors.New("expectedErr")
				fakeEqInterpreter.InterpretReturns(true, expectedError)

				objectUnderTest := _interpreter{
					eqInterpreter: fakeEqInterpreter,
				}

				/* act */
				actualResult, actualError := objectUnderTest.Interpret(
					providedOpHandle,
					providedScgPredicate,
					providedScope,
				)

				/* assert */
				actualExpression,
					actualOpHandle,
					actualScope := fakeEqInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualExpression).To(Equal(providedScgPredicate.Eq))
				Expect(actualOpHandle).To(Equal(providedOpHandle))

				Expect(actualResult).To(Equal(expectedResult))
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("Ne predicate", func() {
			It("should call neInterpreter.Interpret w/ expected args & return result", func() {
				/* arrange */
				providedOpHandle := new(data.FakeHandle)

				providedScgPredicate := &model.SCGPredicate{
					Ne: []interface{}{},
				}

				providedScope := map[string]*model.Value{}

				fakeNeInterpreter := new(ne.FakeInterpreter)
				expectedResult := true
				expectedError := errors.New("expectedErr")
				fakeNeInterpreter.InterpretReturns(true, expectedError)

				objectUnderTest := _interpreter{
					neInterpreter: fakeNeInterpreter,
				}

				/* act */
				actualResult, actualError := objectUnderTest.Interpret(
					providedOpHandle,
					providedScgPredicate,
					providedScope,
				)

				/* assert */
				actualExpression,
					actualOpHandle,
					actualScope := fakeNeInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualExpression).To(Equal(providedScgPredicate.Ne))
				Expect(actualOpHandle).To(Equal(providedOpHandle))

				Expect(actualResult).To(Equal(expectedResult))
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("Unexpected predicate", func() {
			It("should return expected result", func() {
				/* arrange */
				providedScgPredicate := &model.SCGPredicate{}
				objectUnderTest := _interpreter{}

				expectedError := fmt.Errorf("unable to interpret predicate; predicate was unexpected type %+v", providedScgPredicate)

				/* act */
				_, actualError := objectUnderTest.Interpret(
					new(data.FakeHandle),
					providedScgPredicate,
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
