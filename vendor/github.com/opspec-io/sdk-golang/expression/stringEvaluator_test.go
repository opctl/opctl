package expression

import (
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

				objectUnderTest := _stringEvaluator{
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

				objectUnderTest := _stringEvaluator{
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
			It("should call data.CoerceToString w/ expected args", func() {
				/* arrange */
				providedExpression := map[string]interface{}{"dummyName": 2.2}

				fakeData := new(data.Fake)

				objectUnderTest := _stringEvaluator{
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
				Expect(*actualValue).To(Equal(model.Value{Object: providedExpression}))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeData := new(data.Fake)
				coercedValue := model.Value{Object: map[string]interface{}{}}
				coerceToStringErr := errors.New("dummyError")

				fakeData.CoerceToStringReturns(&coercedValue, coerceToStringErr)

				objectUnderTest := _stringEvaluator{
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
		Context("expression is []interface{}", func() {
			It("should call data.CoerceToString w/ expected args", func() {
				/* arrange */
				providedExpression := []interface{}{"dummyName"}

				fakeData := new(data.Fake)

				objectUnderTest := _stringEvaluator{
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

				objectUnderTest := _stringEvaluator{
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

				objectUnderTest := _stringEvaluator{
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

					objectUnderTest := _stringEvaluator{
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

			objectUnderTest := _stringEvaluator{
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
