package expression

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/pkg/errors"
)

var _ = Context("EvalToString", func() {
	var _ = Context("EvalToString", func() {
		Context("expression is float64", func() {
			It("should call data.CoerceToString w/ expected args", func() {
				/* arrange */
				providedExpression := 2.2

				fakeData := new(data.Fake)

				objectUnderTest := _evalStringer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.EvalToString(
					map[string]*model.Value{},
					providedExpression,
					new(pkg.FakeHandle),
				)

				/* assert */
				actualValue := fakeData.CoerceToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Number: &providedExpression}))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeData := new(data.Fake)
				coercedValue := model.Value{Number: new(float64)}
				coerceToStringErr := errors.New("dummyError")

				fakeData.CoerceToStringReturns(&coercedValue, coerceToStringErr)

				objectUnderTest := _evalStringer{
					data: fakeData,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.EvalToString(
					map[string]*model.Value{},
					2.2,
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(*actualValue).To(Equal(coercedValue))
				Expect(actualErr).To(Equal(coerceToStringErr))
			})
		})
		Context("expression is map[string]interface{}", func() {
			It("should call evalObjectInitializerer.Eval w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedExpression := map[string]interface{}{
					"prop1Name": "prop1Value",
				}
				providedPkgRef := new(pkg.FakeHandle)

				fakeEvalObjectInitializerer := new(fakeEvalObjectInitializerer)
				// err to trigger immediate return
				evalErr := errors.New("evalErr")
				fakeEvalObjectInitializerer.EvalReturns(map[string]interface{}{}, evalErr)

				objectUnderTest := _evalStringer{
					evalObjectInitializerer: fakeEvalObjectInitializerer,
				}

				/* act */
				objectUnderTest.EvalToString(
					providedScope,
					providedExpression,
					providedPkgRef,
				)

				/* assert */
				actualExpression,
					actualScope,
					actualPkgRef := fakeEvalObjectInitializerer.EvalArgsForCall(0)

				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualPkgRef).To(Equal(providedPkgRef))

			})
			Context("evalObjectInitializerer.Eval errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := map[string]interface{}{
						"prop1Name": "prop1Value",
					}

					fakeEvalObjectInitializerer := new(fakeEvalObjectInitializerer)
					// err to trigger immediate return
					evalErr := errors.New("evalErr")
					fakeEvalObjectInitializerer.EvalReturns(map[string]interface{}{}, evalErr)

					expectedErr := fmt.Errorf("unable to evaluate %+v to string; error was %v", providedExpression, evalErr)

					objectUnderTest := _evalStringer{
						evalObjectInitializerer: fakeEvalObjectInitializerer,
					}

					/* act */
					_, actualErr := objectUnderTest.EvalToString(
						map[string]*model.Value{},
						providedExpression,
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("evalObjectInitializerer.Eval doesn't err", func() {
				It("should call data.CoerceToString w/ expected args", func() {
					/* arrange */
					expectedObjectValue := map[string]interface{}{"dummyName": 2.2}

					fakeEvalObjectInitializerer := new(fakeEvalObjectInitializerer)
					fakeEvalObjectInitializerer.EvalReturns(expectedObjectValue, nil)

					fakeData := new(data.Fake)

					objectUnderTest := _evalStringer{
						data: fakeData,
						evalObjectInitializerer: fakeEvalObjectInitializerer,
					}

					/* act */
					objectUnderTest.EvalToString(
						map[string]*model.Value{},
						map[string]interface{}{},
						new(pkg.FakeHandle),
					)

					/* assert */
					actualValue := fakeData.CoerceToStringArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Value{Object: expectedObjectValue}))
				})
				It("should return expected result", func() {
					/* arrange */
					fakeData := new(data.Fake)
					coercedValue := model.Value{Object: map[string]interface{}{}}
					coerceToStringErr := errors.New("dummyError")

					fakeData.CoerceToStringReturns(&coercedValue, coerceToStringErr)

					objectUnderTest := _evalStringer{
						evalObjectInitializerer: new(fakeEvalObjectInitializerer),
						data: fakeData,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToString(
						map[string]*model.Value{},
						map[string]interface{}{},
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(*actualValue).To(Equal(coercedValue))
					Expect(actualErr).To(Equal(coerceToStringErr))
				})
			})
		})
		Context("expression is []interface{}", func() {
			It("should call data.CoerceToString w/ expected args", func() {
				/* arrange */
				providedExpression := []interface{}{"dummyName"}

				fakeData := new(data.Fake)

				objectUnderTest := _evalStringer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.EvalToString(
					map[string]*model.Value{},
					providedExpression,
					new(pkg.FakeHandle),
				)

				/* assert */
				actualValue := fakeData.CoerceToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Array: providedExpression}))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeData := new(data.Fake)
				coercedValue := model.Value{Array: []interface{}{}}
				coerceToStringErr := errors.New("dummyError")

				fakeData.CoerceToStringReturns(&coercedValue, coerceToStringErr)

				objectUnderTest := _evalStringer{
					data: fakeData,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.EvalToString(
					map[string]*model.Value{},
					[]interface{}{},
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(*actualValue).To(Equal(coercedValue))
				Expect(actualErr).To(Equal(coerceToStringErr))
			})
		})
		Context("expression is string", func() {
			It("should call interpolater.Interpolate w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedExpression := "dummyExpression"
				providedPkgRef := new(pkg.FakeHandle)

				fakeInterpolater := new(interpolater.Fake)
				// err to trigger immediate return
				fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

				objectUnderTest := _evalStringer{
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.EvalToString(
					providedScope,
					providedExpression,
					providedPkgRef,
				)

				/* assert */
				actualExpression,
					actualScope,
					actualPkgRef := fakeInterpolater.InterpolateArgsForCall(0)

				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualPkgRef).To(Equal(providedPkgRef))

			})
			Context("interpolater.Interpolate errs", func() {
				It("should return expected err", func() {
					/* arrange */
					fakeInterpolater := new(interpolater.Fake)
					interpolateErr := errors.New("dummyError")
					fakeInterpolater.InterpolateReturns("", interpolateErr)

					objectUnderTest := _evalStringer{
						interpolater: fakeInterpolater,
					}

					/* act */
					_, actualErr := objectUnderTest.EvalToString(
						map[string]*model.Value{},
						"dummyExpression",
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(interpolateErr))

				})
			})
		})
		It("should call data.CoerceToString w/ expected args & return result", func() {
			/* arrange */
			fakeInterpolater := new(interpolater.Fake)

			interpolatedValue := "dummyString"
			fakeInterpolater.InterpolateReturns(interpolatedValue, nil)

			fakeData := new(data.Fake)

			coercedValue := model.Value{String: new(string)}
			fakeData.CoerceToStringReturns(&coercedValue, nil)

			objectUnderTest := _evalStringer{
				data:         fakeData,
				interpolater: fakeInterpolater,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.EvalToString(
				map[string]*model.Value{},
				"dummyExpression",
				new(pkg.FakeHandle),
			)

			/* assert */
			Expect(*fakeData.CoerceToStringArgsForCall(0)).To(Equal(model.Value{String: &interpolatedValue}))

			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(BeNil())
		})
	})
})
