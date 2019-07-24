package file

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("Interpret", func() {
	Context("expression is ref", func() {

		It("should call referenceInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedScope := map[string]*types.Value{"dummyName": {}}

			providedExpression := "$(providedExpression)"
			providedOpHandle := new(data.FakeHandle)

			fakeReferenceInterpreter := new(reference.FakeInterpreter)
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

				fakeReferenceInterpreter := new(reference.FakeInterpreter)
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
					map[string]*types.Value{},
					providedExpression,
					new(data.FakeHandle),
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

				referencedValue := new(types.Value)
				fakeReferenceInterpreter := new(reference.FakeInterpreter)
				fakeReferenceInterpreter.InterpretReturns(referencedValue, nil)

				coercedValue := new(types.Value)
				fakeCoerce := new(coerce.Fake)
				fakeCoerce.ToFileReturns(coercedValue, nil)

				objectUnderTest := _interpreter{
					coerce:               fakeCoerce,
					referenceInterpreter: fakeReferenceInterpreter,
				}

				/* act */
				actualResultValue, actualResultErr := objectUnderTest.Interpret(
					map[string]*types.Value{},
					"$(providedExpression)",
					new(data.FakeHandle),
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
		providedScope := map[string]*types.Value{"dummyName": {}}
		providedExpression := map[string]interface{}{
			"prop1Name": "prop1Value",
		}
		providedOpHandle := new(data.FakeHandle)

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
			new(data.FakeHandle),
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

			fakeValueInterpreter := new(value.FakeInterpreter)
			// err to trigger immediate return
			interpretErr := errors.New("interpretErr")
			fakeValueInterpreter.InterpretReturns(types.Value{}, interpretErr)

			expectedErr := fmt.Errorf("unable to interpret %+v to file; error was %v", providedExpression, interpretErr)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			_, actualErr := objectUnderTest.Interpret(
				map[string]*types.Value{},
				providedExpression,
				new(data.FakeHandle),
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

			fakeValueInterpreter := new(value.FakeInterpreter)
			expectedObjectValue := types.Value{String: new(string)}
			fakeValueInterpreter.InterpretReturns(expectedObjectValue, nil)

			fakeCoerce := new(coerce.Fake)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
				coerce:           fakeCoerce,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*types.Value{},
				map[string]interface{}{},
				new(data.FakeHandle),
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
			fakeCoerce := new(coerce.Fake)
			coercedValue := types.Value{Object: new(map[string]interface{})}
			toFileErr := errors.New("dummyError")

			fakeCoerce.ToFileReturns(&coercedValue, toFileErr)

			objectUnderTest := _interpreter{
				valueInterpreter: new(value.FakeInterpreter),
				coerce:           fakeCoerce,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				map[string]*types.Value{},
				map[string]interface{}{},
				new(data.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(toFileErr))
		})
	})
})
