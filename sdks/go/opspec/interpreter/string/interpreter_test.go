package string

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
		providedExpression := map[string]interface{}{
			"prop1Name": "prop1Value",
		}
		providedOpRef := new(data.FakeHandle)

		fakeValueInterpreter := new(value.FakeInterpreter)
		// err to trigger immediate return
		interpretErr := errors.New("interpretErr")
		fakeValueInterpreter.InterpretReturns(types.Value{}, interpretErr)

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
		It("should return expected result", func() {

			/* arrange */
			providedExpression := map[string]interface{}{
				"prop1Name": "prop1Value",
			}

			fakeValueInterpreter := new(value.FakeInterpreter)
			// err to trigger immediate return
			interpretErr := errors.New("interpretErr")
			fakeValueInterpreter.InterpretReturns(types.Value{}, interpretErr)

			expectedErr := fmt.Errorf("unable to interpret %+v to string; error was %v", providedExpression, interpretErr)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			_, actualErr := objectUnderTest.Interpret(
				map[string]*types.Value{},
				providedExpression,
				new(data.FakeHandle),
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Context("valueInterpreter.Interpret doesn't err", func() {
		It("should call coerce.ToString w/ expected args", func() {
			/* arrange */
			expectedObjectValue := types.Value{String: new(string)}

			fakeValueInterpreter := new(value.FakeInterpreter)
			fakeValueInterpreter.InterpretReturns(expectedObjectValue, nil)

			fakeCoerce := new(coerce.Fake)

			objectUnderTest := _interpreter{
				coerce:           fakeCoerce,
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*types.Value{},
				map[string]interface{}{},
				new(data.FakeHandle),
			)

			/* assert */
			actualValue := fakeCoerce.ToStringArgsForCall(0)
			Expect(*actualValue).To(Equal(expectedObjectValue))
		})
		It("should return expected result", func() {
			/* arrange */
			fakeCoerce := new(coerce.Fake)
			coercedValue := types.Value{Object: new(map[string]interface{})}
			toStringErr := errors.New("dummyError")

			fakeCoerce.ToStringReturns(&coercedValue, toStringErr)

			objectUnderTest := _interpreter{
				valueInterpreter: new(value.FakeInterpreter),
				coerce:           fakeCoerce,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				map[string]*types.Value{},
				map[string]interface{}{},
				new(data.FakeHandle),
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(toStringErr))
		})
	})
})
