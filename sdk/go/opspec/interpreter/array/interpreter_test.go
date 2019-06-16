package array

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/data"
	"github.com/opctl/opctl/sdk/go/data/coerce"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/value"
)

var _ = Context("Interpret", func() {
	It("should call valueInterpreter.Interpret w/ expected args", func() {
		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := "dummyExpression"
		providedOpRef := new(data.FakeHandle)

		fakeValueInterpreter := new(value.FakeInterpreter)
		// err to trigger immediate return
		fakeValueInterpreter.InterpretReturns(model.Value{}, errors.New("dummyError"))

		objectUnderTest := _interpreter{
			valueInterpreter: fakeValueInterpreter,
		}

		/* act */
		objectUnderTest.Interpret(
			providedScope,
			providedExpression,
			providedOpRef,
		)

		/* assert */
		actualExpression,
			actualScope,
			actualOpRef := fakeValueInterpreter.InterpretArgsForCall(0)

		Expect(actualExpression).To(Equal(providedExpression))
		Expect(actualScope).To(Equal(providedScope))
		Expect(actualOpRef).To(Equal(providedOpRef))

	})
	Context("valueInterpreter.Interpret errs", func() {
		It("should return expected err", func() {
			/* arrange */
			providedExpression := "dummyExpression"

			fakeValueInterpreter := new(value.FakeInterpreter)
			interpretErr := errors.New("dummyError")
			fakeValueInterpreter.InterpretReturns(model.Value{}, interpretErr)

			expectedErr := fmt.Errorf("unable to interpret %+v to array; error was %v", providedExpression, interpretErr)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			_, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				"dummyExpression",
				new(data.FakeHandle),
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))

		})
	})
	Context("valueInterpreter.Interpret doesn't err", func() {
		It("should call coerce.ToArray w/ expected args & return result", func() {
			/* arrange */
			fakeValueInterpreter := new(value.FakeInterpreter)

			expectedValue := model.Value{String: new(string)}
			fakeValueInterpreter.InterpretReturns(expectedValue, nil)

			fakeCoerce := new(coerce.Fake)

			coercedValue := model.Value{Array: new([]interface{})}
			fakeCoerce.ToArrayReturns(&coercedValue, nil)

			objectUnderTest := _interpreter{
				coerce:           fakeCoerce,
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			actualArray, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				"dummyExpression",
				new(data.FakeHandle),
			)

			/* assert */
			Expect(*fakeCoerce.ToArrayArgsForCall(0)).To(Equal(expectedValue))

			Expect(*actualArray).To(Equal(coercedValue))
			Expect(actualErr).To(BeNil())
		})
	})
})
