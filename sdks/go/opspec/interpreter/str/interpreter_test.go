package str

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
		providedExpression := map[string]interface{}{
			"prop1Name": "prop1Value",
		}
		providedOpRef := new(modelFakes.FakeDataHandle)

		fakeValueInterpreter := new(valueFakes.FakeInterpreter)
		// err to trigger immediate return
		interpretErr := errors.New("interpretErr")
		fakeValueInterpreter.InterpretReturns(model.Value{}, interpretErr)

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

			fakeValueInterpreter := new(valueFakes.FakeInterpreter)
			// err to trigger immediate return
			interpretErr := errors.New("interpretErr")
			fakeValueInterpreter.InterpretReturns(model.Value{}, interpretErr)

			expectedErr := fmt.Errorf("unable to interpret %+v to string; error was %v", providedExpression, interpretErr)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			_, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedExpression,
				new(modelFakes.FakeDataHandle),
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Context("valueInterpreter.Interpret doesn't err", func() {
		It("should call coerce.ToString w/ expected args", func() {
			/* arrange */
			expectedObjectValue := model.Value{String: new(string)}

			fakeValueInterpreter := new(valueFakes.FakeInterpreter)
			fakeValueInterpreter.InterpretReturns(expectedObjectValue, nil)

			fakeCoerce := new(coerceFakes.FakeCoerce)

			objectUnderTest := _interpreter{
				coerce:           fakeCoerce,
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				map[string]interface{}{},
				new(modelFakes.FakeDataHandle),
			)

			/* assert */
			actualValue := fakeCoerce.ToStringArgsForCall(0)
			Expect(*actualValue).To(Equal(expectedObjectValue))
		})
		It("should return expected result", func() {
			/* arrange */
			fakeCoerce := new(coerceFakes.FakeCoerce)
			coercedValue := model.Value{Object: new(map[string]interface{})}
			toStringErr := errors.New("dummyError")

			fakeCoerce.ToStringReturns(&coercedValue, toStringErr)

			objectUnderTest := _interpreter{
				valueInterpreter: new(valueFakes.FakeInterpreter),
				coerce:           fakeCoerce,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				map[string]interface{}{},
				new(modelFakes.FakeDataHandle),
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(toStringErr))
		})
	})
})
