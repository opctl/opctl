package array

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/array/initializer"
	"github.com/opspec-io/sdk-golang/op/interpreter/interpolater"
	"github.com/pkg/errors"
)

var _ = Context("Interpret", func() {
	var _ = Context("Interpret", func() {
		Context("expression is []interface{}", func() {

			It("should call initializerInterpreter.Interpret w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedExpression := []interface{}{
					"item1",
				}
				providedOpRef := new(data.FakeHandle)

				fakeInitializerInterpreter := new(initializer.FakeInterpreter)
				// err to trigger immediate return
				evalErr := errors.New("evalErr")
				fakeInitializerInterpreter.InterpretReturns([]interface{}{}, evalErr)

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
					providedExpression := []interface{}{
						"item1",
					}

					fakeInitializerInterpreter := new(initializer.FakeInterpreter)
					// err to trigger immediate return
					evalErr := errors.New("evalErr")
					fakeInitializerInterpreter.InterpretReturns([]interface{}{}, evalErr)

					expectedErr := fmt.Errorf("unable to evaluate %+v to array; error was %v", providedExpression, evalErr)

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
					expectedResult := []interface{}{"arrayItem"}
					fakeInitializerInterpreter.InterpretReturns(expectedResult, nil)

					objectUnderTest := _interpreter{
						initializerInterpreter: fakeInitializerInterpreter,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						[]interface{}{},
						new(data.FakeHandle),
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Array: expectedResult}))
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
				It("should call coerce.ToArray w/ expected args & return result", func() {
					/* arrange */
					fakeInterpolater := new(interpolater.Fake)

					interpolatedValue := "dummyString"
					fakeInterpolater.InterpolateReturns(interpolatedValue, nil)

					fakeCoerce := new(coerce.Fake)

					coercedValue := model.Value{Array: []interface{}{"arrayItem"}}
					fakeCoerce.ToArrayReturns(&coercedValue, nil)

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
					Expect(*fakeCoerce.ToArrayArgsForCall(0)).To(Equal(model.Value{String: &interpolatedValue}))

					Expect(*actualValue).To(Equal(coercedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
