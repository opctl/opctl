package array

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coerceFakes "github.com/opctl/opctl/sdks/go/data/coerce/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	valueFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/value/fakes"
)

var _ = Context("Interpret", func() {
	It("should call valueInterpreter.Interpret w/ expected args", func() {
		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := "dummyExpression"
		providedOpRef := new(modelFakes.FakeDataHandle)

		fakeValueInterpreter := new(valueFakes.FakeInterpreter)
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

			fakeValueInterpreter := new(valueFakes.FakeInterpreter)
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
				new(modelFakes.FakeDataHandle),
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))

		})
	})
	Context("valueInterpreter.Interpret doesn't err", func() {
		It("should call coerce.ToArray w/ expected args & return result", func() {
			/* arrange */
			fakeValueInterpreter := new(valueFakes.FakeInterpreter)

			expectedValue := model.Value{String: new(string)}
			fakeValueInterpreter.InterpretReturns(expectedValue, nil)

			fakeCoerce := new(coerceFakes.FakeCoerce)

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
				new(modelFakes.FakeDataHandle),
			)

			/* assert */
			Expect(*fakeCoerce.ToArrayArgsForCall(0)).To(Equal(expectedValue))

			Expect(*actualArray).To(Equal(coercedValue))
			Expect(actualErr).To(BeNil())
		})
	})
})
