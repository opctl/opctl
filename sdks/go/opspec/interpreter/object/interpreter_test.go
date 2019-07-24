package object

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("Interpret", func() {
	It("should call valueInterpreter.Interpret w/ expected args", func() {
		/* arrange */
		providedScope := map[string]*types.Value{"dummyName": {}}
		providedExpression := "dummyExpression"
		providedOpRef := new(data.FakeHandle)

		fakeValueInterpreter := new(value.FakeInterpreter)
		// err to trigger immediate return
		fakeValueInterpreter.InterpretReturns(types.Value{}, errors.New("dummyError"))

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
			fakeValueInterpreter.InterpretReturns(types.Value{}, interpretErr)

			expectedErr := fmt.Errorf("unable to interpret %+v to object; error was %v", providedExpression, interpretErr)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			_, actualErr := objectUnderTest.Interpret(
				map[string]*types.Value{},
				"dummyExpression",
				new(data.FakeHandle),
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))

		})
	})
	Context("valueInterpreter.Interpret doesn't err", func() {
		It("should call coerce.ToObject w/ expected args & return result", func() {
			/* arrange */
			fakeValueInterpreter := new(value.FakeInterpreter)

			expectedValue := types.Value{String: new(string)}
			fakeValueInterpreter.InterpretReturns(expectedValue, nil)

			fakeCoerce := new(coerce.Fake)

			coercedValue := types.Value{Object: new(map[string]interface{})}
			fakeCoerce.ToObjectReturns(&coercedValue, nil)

			objectUnderTest := _interpreter{
				coerce:           fakeCoerce,
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			actualObject, actualErr := objectUnderTest.Interpret(
				map[string]*types.Value{},
				"dummyExpression",
				new(data.FakeHandle),
			)

			/* assert */
			Expect(*fakeCoerce.ToObjectArgsForCall(0)).To(Equal(expectedValue))

			Expect(*actualObject).To(Equal(coercedValue))
			Expect(actualErr).To(BeNil())
		})
	})
})
