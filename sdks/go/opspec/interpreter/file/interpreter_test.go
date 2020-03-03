package file

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coerceFakes "github.com/opctl/opctl/sdks/go/data/coerce/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	referenceFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/fakes"
	valueFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/value/fakes"
)

var _ = Context("Interpret", func() {
	Context("expression is ref", func() {

		It("should call referenceInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedScope := map[string]*model.Value{"dummyName": {}}

			providedExpression := "$(providedExpression)"
			providedOpHandle := new(modelFakes.FakeDataHandle)

			fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
			// err to trigger immediate return
			fakeReferenceInterpreter.InterpretReturns(nil, errors.New("dummyError"))

			objectUnderTest := _interpreter{
				referenceInterpreter: fakeReferenceInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedExpression,
				providedOpHandle,
				"dummyScratchDir",
			)

			/* assert */
			actualExpression,
				actualScope,
				actualOpHandle := fakeReferenceInterpreter.InterpretArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
		})
		Context("referenceInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedExpression := "$(providedExpression)"

				fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
				interpretErr := errors.New("dummyError")
				fakeReferenceInterpreter.InterpretReturns(nil, errors.New("dummyError"))

				expectedErr := fmt.Errorf(
					"unable to interpret %+v to file; error was %v",
					providedExpression,
					interpretErr,
				)

				objectUnderTest := _interpreter{
					referenceInterpreter: fakeReferenceInterpreter,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedExpression,
					new(modelFakes.FakeDataHandle),
					"providedScratchDir",
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("referenceInterpreter.Interpret doesn't error", func() {
			It("should call coerce.ToFile w/ expected args & return result", func() {
				/* arrange */
				providedScratchDir := "providedScratchDir"

				referencedValue := new(model.Value)
				fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
				fakeReferenceInterpreter.InterpretReturns(referencedValue, nil)

				coercedValue := new(model.Value)
				fakeCoerce := new(coerceFakes.FakeCoerce)
				fakeCoerce.ToFileReturns(coercedValue, nil)

				objectUnderTest := _interpreter{
					coerce:               fakeCoerce,
					referenceInterpreter: fakeReferenceInterpreter,
				}

				/* act */
				actualResultValue, actualResultErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					"$(providedExpression)",
					new(modelFakes.FakeDataHandle),
					providedScratchDir,
				)

				/* assert */
				Expect(actualResultValue).To(Equal(coercedValue))
				Expect(actualResultErr).To(BeNil())

				actualValue,
					actualScratchDir := fakeCoerce.ToFileArgsForCall(0)
				Expect(actualValue).To(Equal(referencedValue))
				Expect(actualScratchDir).To(Equal(providedScratchDir))
			})
		})
	})
	It("should call valueInterpreter.Interpret w/ expected args", func() {

		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := map[string]interface{}{
			"prop1Name": "prop1Value",
		}
		providedOpHandle := new(modelFakes.FakeDataHandle)

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
			new(modelFakes.FakeDataHandle),
			"dummyScratchDir",
		)

		/* assert */
		actualExpression,
			actualScope,
			actualOpHandle := fakeValueInterpreter.InterpretArgsForCall(0)

		Expect(actualExpression).To(Equal(providedExpression))
		Expect(actualScope).To(Equal(providedScope))
		Expect(actualOpHandle).To(Equal(providedOpHandle))

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

			expectedErr := fmt.Errorf("unable to interpret %+v to file; error was %v", providedExpression, interpretErr)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			_, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedExpression,
				new(modelFakes.FakeDataHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Context("valueInterpreter.Interpret doesn't err", func() {
		It("should call coerce.ToFile w/ expected args", func() {
			/* arrange */
			providedScratchDir := "dummyScratchDir"

			fakeValueInterpreter := new(valueFakes.FakeInterpreter)
			expectedObjectValue := model.Value{String: new(string)}
			fakeValueInterpreter.InterpretReturns(expectedObjectValue, nil)

			fakeCoerce := new(coerceFakes.FakeCoerce)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
				coerce:           fakeCoerce,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				map[string]interface{}{},
				new(modelFakes.FakeDataHandle),
				providedScratchDir,
			)

			/* assert */
			actualValue,
				actualScratchDir := fakeCoerce.ToFileArgsForCall(0)
			Expect(*actualValue).To(Equal(expectedObjectValue))
			Expect(actualScratchDir).To(Equal(providedScratchDir))
		})
		It("should return expected result", func() {
			/* arrange */
			fakeCoerce := new(coerceFakes.FakeCoerce)
			coercedValue := model.Value{Object: new(map[string]interface{})}
			toFileErr := errors.New("dummyError")

			fakeCoerce.ToFileReturns(&coercedValue, toFileErr)

			objectUnderTest := _interpreter{
				valueInterpreter: new(valueFakes.FakeInterpreter),
				coerce:           fakeCoerce,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				map[string]interface{}{},
				new(modelFakes.FakeDataHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(toFileErr))
		})
	})
})
