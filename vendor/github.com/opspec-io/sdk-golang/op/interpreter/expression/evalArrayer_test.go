package expression

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
	"github.com/pkg/errors"
)

var _ = Context("EvalToArray", func() {
	var _ = Context("EvalToArray", func() {
		Context("expression is []interface{}", func() {

			It("should call evalArrayInitializer.Eval w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedExpression := []interface{}{
					"item1",
				}
				providedPkgRef := new(data.FakeHandle)

				fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
				// err to trigger immediate return
				evalErr := errors.New("evalErr")
				fakeEvalArrayInitializerer.EvalReturns([]interface{}{}, evalErr)

				arrayUnderTest := _evalArrayer{
					evalArrayInitializerer: fakeEvalArrayInitializerer,
				}

				/* act */
				arrayUnderTest.EvalToArray(
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
			Context("evalArrayInitializer.Eval errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := []interface{}{
						"item1",
					}

					fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
					// err to trigger immediate return
					evalErr := errors.New("evalErr")
					fakeEvalArrayInitializerer.EvalReturns([]interface{}{}, evalErr)

					expectedErr := fmt.Errorf("unable to evaluate %+v to array; error was %v", providedExpression, evalErr)

					arrayUnderTest := _evalArrayer{
						evalArrayInitializerer: fakeEvalArrayInitializerer,
					}

					/* act */
					_, actualErr := arrayUnderTest.EvalToArray(
						map[string]*model.Value{},
						providedExpression,
						new(data.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))

				})

			})
			Context("evalArrayInitializer.Eval doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
					expectedResult := []interface{}{"arrayItem"}
					fakeEvalArrayInitializerer.EvalReturns(expectedResult, nil)

					arrayUnderTest := _evalArrayer{
						evalArrayInitializerer: fakeEvalArrayInitializerer,
					}

					/* act */
					actualValue, actualErr := arrayUnderTest.EvalToArray(
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
				providedPkgRef := new(data.FakeHandle)

				fakeInterpolater := new(interpolater.Fake)
				// err to trigger immediate return
				fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

				arrayUnderTest := _evalArrayer{
					interpolater: fakeInterpolater,
				}

				/* act */
				arrayUnderTest.EvalToArray(
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

					arrayUnderTest := _evalArrayer{
						interpolater: fakeInterpolater,
					}

					/* act */
					_, actualErr := arrayUnderTest.EvalToArray(
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

					arrayUnderTest := _evalArrayer{
						coerce:       fakeCoerce,
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := arrayUnderTest.EvalToArray(
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
