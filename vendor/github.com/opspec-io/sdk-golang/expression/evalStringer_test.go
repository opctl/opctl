package expression

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/pkg/errors"
)

var _ = Context("EvalToString", func() {
	var _ = Context("EvalToString", func() {
		Context("expression is float64", func() {
			It("should call coerce.ToString w/ expected args", func() {
				/* arrange */
				providedExpression := 2.2

				fakeCoerce := new(coerce.Fake)

				objectUnderTest := _evalStringer{
					coerce: fakeCoerce,
				}

				/* act */
				objectUnderTest.EvalToString(
					map[string]*model.Value{},
					providedExpression,
					new(pkg.FakeHandle),
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

				objectUnderTest := _evalStringer{
					coerce: fakeCoerce,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.EvalToString(
					map[string]*model.Value{},
					2.2,
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(*actualValue).To(Equal(coercedValue))
				Expect(actualErr).To(Equal(toStringErr))
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
				It("should call coerce.ToString w/ expected args", func() {
					/* arrange */
					expectedObjectValue := map[string]interface{}{"dummyName": 2.2}

					fakeEvalObjectInitializerer := new(fakeEvalObjectInitializerer)
					fakeEvalObjectInitializerer.EvalReturns(expectedObjectValue, nil)

					fakeCoerce := new(coerce.Fake)

					objectUnderTest := _evalStringer{
						coerce:                  fakeCoerce,
						evalObjectInitializerer: fakeEvalObjectInitializerer,
					}

					/* act */
					objectUnderTest.EvalToString(
						map[string]*model.Value{},
						map[string]interface{}{},
						new(pkg.FakeHandle),
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

					objectUnderTest := _evalStringer{
						evalObjectInitializerer: new(fakeEvalObjectInitializerer),
						coerce:                  fakeCoerce,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToString(
						map[string]*model.Value{},
						map[string]interface{}{},
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(*actualValue).To(Equal(coercedValue))
					Expect(actualErr).To(Equal(toStringErr))
				})
			})
		})
		Context("expression is []interface{}", func() {
			It("should call evalArrayInitializerer.Eval w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedExpression := []interface{}{
					"item1",
				}
				providedPkgRef := new(pkg.FakeHandle)

				fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
				// err to trigger immediate return
				evalErr := errors.New("evalErr")
				fakeEvalArrayInitializerer.EvalReturns([]interface{}{}, evalErr)

				arrayUnderTest := _evalStringer{
					evalArrayInitializerer: fakeEvalArrayInitializerer,
				}

				/* act */
				arrayUnderTest.EvalToString(
					providedScope,
					providedExpression,
					providedPkgRef,
				)

				/* assert */
				actualExpression,
					actualScope,
					actualPkgRef := fakeEvalArrayInitializerer.EvalArgsForCall(0)

				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualPkgRef).To(Equal(providedPkgRef))

			})
			Context("evalArrayInitializerer.Eval errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := []interface{}{
						"item1",
					}

					fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
					// err to trigger immediate return
					evalErr := errors.New("evalErr")
					fakeEvalArrayInitializerer.EvalReturns([]interface{}{}, evalErr)

					expectedErr := fmt.Errorf("unable to evaluate %+v to string; error was %v", providedExpression, evalErr)

					arrayUnderTest := _evalStringer{
						evalArrayInitializerer: fakeEvalArrayInitializerer,
					}

					/* act */
					_, actualErr := arrayUnderTest.EvalToString(
						map[string]*model.Value{},
						providedExpression,
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("evalArrayInitializerer.Eval doesn't err", func() {
				It("should call coerce.ToString w/ expected args", func() {
					/* arrange */
					expectedArrayValue := []interface{}{"item1"}

					fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
					fakeEvalArrayInitializerer.EvalReturns(expectedArrayValue, nil)

					fakeCoerce := new(coerce.Fake)

					arrayUnderTest := _evalStringer{
						coerce:                 fakeCoerce,
						evalArrayInitializerer: fakeEvalArrayInitializerer,
					}

					/* act */
					arrayUnderTest.EvalToString(
						map[string]*model.Value{},
						[]interface{}{},
						new(pkg.FakeHandle),
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

					arrayUnderTest := _evalStringer{
						evalArrayInitializerer: new(fakeEvalArrayInitializerer),
						coerce:                 fakeCoerce,
					}

					/* act */
					actualValue, actualErr := arrayUnderTest.EvalToString(
						map[string]*model.Value{},
						[]interface{}{},
						new(pkg.FakeHandle),
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
		It("should call coerce.ToString w/ expected args & return result", func() {
			/* arrange */
			fakeInterpolater := new(interpolater.Fake)

			interpolatedValue := "dummyString"
			fakeInterpolater.InterpolateReturns(interpolatedValue, nil)

			fakeCoerce := new(coerce.Fake)

			coercedValue := model.Value{String: new(string)}
			fakeCoerce.ToStringReturns(&coercedValue, nil)

			objectUnderTest := _evalStringer{
				coerce:       fakeCoerce,
				interpolater: fakeInterpolater,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.EvalToString(
				map[string]*model.Value{},
				"dummyExpression",
				new(pkg.FakeHandle),
			)

			/* assert */
			Expect(*fakeCoerce.ToStringArgsForCall(0)).To(Equal(model.Value{String: &interpolatedValue}))

			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(BeNil())
		})
	})
})
