package string

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	arrayInitializer "github.com/opctl/sdk-golang/opspec/interpreter/array/initializer"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
	objectInitializer "github.com/opctl/sdk-golang/opspec/interpreter/object/initializer"
)

var _ = Context("Interpret", func() {
	var _ = Context("Interpret", func() {
		Context("expression is bool", func() {
			It("should call coerce.ToString w/ expected args", func() {
				/* arrange */
				providedExpression := false

				fakeCoerce := new(coerce.Fake)

				objectUnderTest := _interpreter{
					coerce: fakeCoerce,
				}

				/* act */
				objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedExpression,
					new(data.FakeHandle),
				)

				/* assert */
				actualValue := fakeCoerce.ToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Boolean: &providedExpression}))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeCoerce := new(coerce.Fake)
				coercedValue := model.Value{Boolean: new(bool)}
				toStringErr := errors.New("dummyError")

				fakeCoerce.ToStringReturns(&coercedValue, toStringErr)

				objectUnderTest := _interpreter{
					coerce: fakeCoerce,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					false,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(*actualValue).To(Equal(coercedValue))
				Expect(actualErr).To(Equal(toStringErr))
			})
		})
		Context("expression is float64", func() {
			It("should call coerce.ToString w/ expected args", func() {
				/* arrange */
				providedExpression := 2.2

				fakeCoerce := new(coerce.Fake)

				objectUnderTest := _interpreter{
					coerce: fakeCoerce,
				}

				/* act */
				objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedExpression,
					new(data.FakeHandle),
				)

				/* assert */
				actualValue := fakeCoerce.ToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Number: &providedExpression}))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeCoerce := new(coerce.Fake)
				coercedValue := model.Value{Number: new(float64)}
				toStringErr := errors.New("dummyError")

				fakeCoerce.ToStringReturns(&coercedValue, toStringErr)

				objectUnderTest := _interpreter{
					coerce: fakeCoerce,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					2.2,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(*actualValue).To(Equal(coercedValue))
				Expect(actualErr).To(Equal(toStringErr))
			})
		})
		Context("expression is int", func() {
			It("should call coerce.ToString w/ expected args", func() {
				/* arrange */
				providedExpression := 2

				expectedNumber := float64(providedExpression)

				fakeCoerce := new(coerce.Fake)

				objectUnderTest := _interpreter{
					coerce: fakeCoerce,
				}

				/* act */
				objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedExpression,
					new(data.FakeHandle),
				)

				/* assert */
				actualValue := fakeCoerce.ToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Number: &expectedNumber}))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeCoerce := new(coerce.Fake)
				coercedValue := model.Value{Number: new(float64)}
				toStringErr := errors.New("dummyError")

				fakeCoerce.ToStringReturns(&coercedValue, toStringErr)

				objectUnderTest := _interpreter{
					coerce: fakeCoerce,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					2,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(*actualValue).To(Equal(coercedValue))
				Expect(actualErr).To(Equal(toStringErr))
			})
		})
		Context("expression is map[string]interface{}", func() {
			It("should call objectInitializerInterpreter.Interpret w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedExpression := map[string]interface{}{
					"prop1Name": "prop1Value",
				}
				providedOpRef := new(data.FakeHandle)

				fakeObjectInitializerInterpreter := new(objectInitializer.FakeInterpreter)
				// err to trigger immediate return
				interpretErr := errors.New("interpretErr")
				fakeObjectInitializerInterpreter.InterpretReturns(map[string]interface{}{}, interpretErr)

				objectUnderTest := _interpreter{
					objectInitializerInterpreter: fakeObjectInitializerInterpreter,
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
					actualOpRef := fakeObjectInitializerInterpreter.InterpretArgsForCall(0)

				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpRef).To(Equal(providedOpRef))

			})
			Context("interpreter.Eval errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := map[string]interface{}{
						"prop1Name": "prop1Value",
					}

					fakeObjectInitializerInterpreter := new(objectInitializer.FakeInterpreter)
					// err to trigger immediate return
					interpretErr := errors.New("interpretErr")
					fakeObjectInitializerInterpreter.InterpretReturns(map[string]interface{}{}, interpretErr)

					expectedErr := fmt.Errorf("unable to evaluate %+v to string; error was %v", providedExpression, interpretErr)

					objectUnderTest := _interpreter{
						objectInitializerInterpreter: fakeObjectInitializerInterpreter,
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
			Context("interpreter.Eval doesn't err", func() {
				It("should call coerce.ToString w/ expected args", func() {
					/* arrange */
					expectedObjectValue := map[string]interface{}{"dummyName": 2.2}

					fakeObjectInitializerInterpreter := new(objectInitializer.FakeInterpreter)
					fakeObjectInitializerInterpreter.InterpretReturns(expectedObjectValue, nil)

					fakeCoerce := new(coerce.Fake)

					objectUnderTest := _interpreter{
						coerce: fakeCoerce,
						objectInitializerInterpreter: fakeObjectInitializerInterpreter,
					}

					/* act */
					objectUnderTest.Interpret(
						map[string]*model.Value{},
						map[string]interface{}{},
						new(data.FakeHandle),
					)

					/* assert */
					actualValue := fakeCoerce.ToStringArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Value{Object: expectedObjectValue}))
				})
				It("should return expected result", func() {
					/* arrange */
					fakeCoerce := new(coerce.Fake)
					coercedValue := model.Value{Object: map[string]interface{}{}}
					toStringErr := errors.New("dummyError")

					fakeCoerce.ToStringReturns(&coercedValue, toStringErr)

					objectUnderTest := _interpreter{
						objectInitializerInterpreter: new(objectInitializer.FakeInterpreter),
						coerce: fakeCoerce,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						map[string]interface{}{},
						new(data.FakeHandle),
					)

					/* assert */
					Expect(*actualValue).To(Equal(coercedValue))
					Expect(actualErr).To(Equal(toStringErr))
				})
			})
		})
		Context("expression is []interface{}", func() {
			It("should call arrayInitializerInterpreter.Interpret w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedExpression := []interface{}{
					"item1",
				}
				providedOpRef := new(data.FakeHandle)

				fakeArrayInitializerInterpreter := new(arrayInitializer.FakeInterpreter)
				// err to trigger immediate return
				interpretErr := errors.New("interpretErr")
				fakeArrayInitializerInterpreter.InterpretReturns([]interface{}{}, interpretErr)

				objectUnderTest := _interpreter{
					arrayInitializerInterpreter: fakeArrayInitializerInterpreter,
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
					actualOpRef := fakeArrayInitializerInterpreter.InterpretArgsForCall(0)

				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpRef).To(Equal(providedOpRef))

			})
			Context("arrayInitializerInterpreter.Interpret errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := []interface{}{
						"item1",
					}

					fakeArrayInitializerInterpreter := new(arrayInitializer.FakeInterpreter)
					// err to trigger immediate return
					interpretErr := errors.New("interpretErr")
					fakeArrayInitializerInterpreter.InterpretReturns([]interface{}{}, interpretErr)

					expectedErr := fmt.Errorf("unable to evaluate %+v to string; error was %v", providedExpression, interpretErr)

					objectUnderTest := _interpreter{
						arrayInitializerInterpreter: fakeArrayInitializerInterpreter,
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
			Context("arrayInitializerInterpreter.Interpret doesn't err", func() {
				It("should call coerce.ToString w/ expected args", func() {
					/* arrange */
					expectedArrayValue := []interface{}{"item1"}

					fakeArrayInitializerInterpreter := new(arrayInitializer.FakeInterpreter)
					fakeArrayInitializerInterpreter.InterpretReturns(expectedArrayValue, nil)

					fakeCoerce := new(coerce.Fake)

					objectUnderTest := _interpreter{
						coerce: fakeCoerce,
						arrayInitializerInterpreter: fakeArrayInitializerInterpreter,
					}

					/* act */
					objectUnderTest.Interpret(
						map[string]*model.Value{},
						[]interface{}{},
						new(data.FakeHandle),
					)

					/* assert */
					actualValue := fakeCoerce.ToStringArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Value{Array: expectedArrayValue}))
				})
				It("should return expected result", func() {
					/* arrange */
					fakeCoerce := new(coerce.Fake)
					coercedValue := model.Value{Array: []interface{}{}}
					toStringErr := errors.New("dummyError")

					fakeCoerce.ToStringReturns(&coercedValue, toStringErr)

					objectUnderTest := _interpreter{
						arrayInitializerInterpreter: new(arrayInitializer.FakeInterpreter),
						coerce: fakeCoerce,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						[]interface{}{},
						new(data.FakeHandle),
					)

					/* assert */
					Expect(*actualValue).To(Equal(coercedValue))
					Expect(actualErr).To(Equal(toStringErr))
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
		})
		It("should call coerce.ToString w/ expected args & return result", func() {
			/* arrange */
			fakeInterpolater := new(interpolater.Fake)

			interpolatedValue := "dummyString"
			fakeInterpolater.InterpolateReturns(interpolatedValue, nil)

			fakeCoerce := new(coerce.Fake)

			coercedValue := model.Value{String: new(string)}
			fakeCoerce.ToStringReturns(&coercedValue, nil)

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
			Expect(*fakeCoerce.ToStringArgsForCall(0)).To(Equal(model.Value{String: &interpolatedValue}))

			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(BeNil())
		})
	})
})
