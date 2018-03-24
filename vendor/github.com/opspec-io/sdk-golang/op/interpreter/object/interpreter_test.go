package object

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/interpolater"
	"github.com/opspec-io/sdk-golang/op/interpreter/object/initializer"
)

var _ = Context("Interpret", func() {
	Context("expression is map[string]interface{}", func() {
		It("should call initializerInterpreter.Interpret w/ expected args", func() {

			/* arrange */
			providedScope := map[string]*model.Value{"dummyName": {}}
			providedExpression := map[string]interface{}{
				"prop1Name": "prop1Value",
			}
			providedOpRef := new(data.FakeHandle)

			fakeInitializerInterpreter := new(initializer.FakeInterpreter)
			// err to trigger immediate return
			interpretErr := errors.New("interpretErr")
			fakeInitializerInterpreter.InterpretReturns(map[string]interface{}{}, interpretErr)

			objectUnderTest := _interpreter{
				initializerInterpreter: fakeInitializerInterpreter,
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
				actualOpRef := fakeInitializerInterpreter.InterpretArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpRef).To(Equal(providedOpRef))

		})
		Context("initializerInterpreter.Interpret errs", func() {
			It("should return expected result", func() {

				/* arrange */
				providedExpression := map[string]interface{}{
					"prop1Name": "prop1Value",
				}

				fakeInitializerInterpreter := new(initializer.FakeInterpreter)
				// err to trigger immediate return
				interpretErr := errors.New("interpretErr")
				fakeInitializerInterpreter.InterpretReturns(map[string]interface{}{}, interpretErr)

				expectedErr := fmt.Errorf("unable to interpretuate %+v to object; error was %v", providedExpression, interpretErr)

				objectUnderTest := _interpreter{
					initializerInterpreter: fakeInitializerInterpreter,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedExpression,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))

			})

		})
		Context("initializerInterpreter.Interpret doesn't err", func() {
			It("should return expected result", func() {

				/* arrange */
				fakeInitializerInterpreter := new(initializer.FakeInterpreter)
				expectedValue := map[string]interface{}{
					"prop1Name": "prop1Value",
				}
				fakeInitializerInterpreter.InterpretReturns(expectedValue, nil)

				objectUnderTest := _interpreter{
					initializerInterpreter: fakeInitializerInterpreter,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					map[string]interface{}{},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Object: expectedValue}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("expression is string", func() {
		It("should call interpolater.Interpolate w/ expected args", func() {
			/* arrange */
			providedScope := map[string]*model.Value{"dummyName": {}}
			providedExpression := "dummyExpression"
			providedOpRef := new(data.FakeHandle)

			fakeInterpolater := new(interpolater.Fake)
			// err to trigger immediate return
			fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

			objectUnderTest := _interpreter{
				interpolater: fakeInterpolater,
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
				actualOpRef := fakeInterpolater.InterpolateArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpRef).To(Equal(providedOpRef))

		})
		Context("interpolater.Interpolate errs", func() {
			It("should return expected err", func() {
				/* arrange */
				fakeInterpolater := new(interpolater.Fake)
				interpolateErr := errors.New("dummyError")
				fakeInterpolater.InterpolateReturns("", interpolateErr)

				objectUnderTest := _interpreter{
					interpolater: fakeInterpolater,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					"dummyExpression",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(interpolateErr))

			})
		})
		Context("interpolater.Interpolate doesn't err", func() {
			It("should call coerce.ToObject w/ expected args & return result", func() {
				/* arrange */
				fakeInterpolater := new(interpolater.Fake)

				interpolatedValue := "dummyString"
				fakeInterpolater.InterpolateReturns(interpolatedValue, nil)

				fakeCoerce := new(coerce.Fake)

				coercedValue := model.Value{Object: map[string]interface{}{"dummyName": "dummyValue"}}
				fakeCoerce.ToObjectReturns(&coercedValue, nil)

				objectUnderTest := _interpreter{
					coerce:       fakeCoerce,
					interpolater: fakeInterpolater,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					"dummyExpression",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(*fakeCoerce.ToObjectArgsForCall(0)).To(Equal(model.Value{String: &interpolatedValue}))

				Expect(*actualValue).To(Equal(coercedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
