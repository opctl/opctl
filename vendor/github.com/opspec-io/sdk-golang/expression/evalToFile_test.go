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

var _ = Context("EvalToFile", func() {
	Context("expression is float64", func() {
		It("should call data.CoerceToFile w/ expected args", func() {
			/* arrange */
			providedExpression := 2.2
			providedScratchDir := "dummyScratchDir"

			fakeData := new(data.Fake)

			objectUnderTest := _evalToFile{
				data: fakeData,
			}

			/* act */
			objectUnderTest.EvalToFile(
				map[string]*model.Value{},
				providedExpression,
				new(pkg.FakeHandle),
				providedScratchDir,
			)

			/* assert */
			actualValue,
				actualScratchDir := fakeData.CoerceToFileArgsForCall(0)
			Expect(*actualValue).To(Equal(model.Value{Number: &providedExpression}))
			Expect(actualScratchDir).To(Equal(providedScratchDir))
		})
		It("should return expected result", func() {
			/* arrange */
			fakeData := new(data.Fake)
			coercedValue := model.Value{Number: new(float64)}
			coerceToFileErr := errors.New("dummyError")

			fakeData.CoerceToFileReturns(&coercedValue, coerceToFileErr)

			objectUnderTest := _evalToFile{
				data: fakeData,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.EvalToFile(
				map[string]*model.Value{},
				2.2,
				new(pkg.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(coerceToFileErr))
		})
	})
	Context("expression is map[string]interface{}", func() {
		It("should call data.CoerceToFile w/ expected args", func() {
			/* arrange */
			providedExpression := map[string]interface{}{"dummyName": 2.2}
			providedScratchDir := "dummyScratchDir"

			fakeData := new(data.Fake)

			objectUnderTest := _evalToFile{
				data: fakeData,
			}

			/* act */
			objectUnderTest.EvalToFile(
				map[string]*model.Value{},
				providedExpression,
				new(pkg.FakeHandle),
				providedScratchDir,
			)

			/* assert */
			actualValue,
				actualScratchDir := fakeData.CoerceToFileArgsForCall(0)
			Expect(*actualValue).To(Equal(model.Value{Object: providedExpression}))
			Expect(actualScratchDir).To(Equal(providedScratchDir))
		})
		It("should return expected result", func() {
			/* arrange */
			fakeData := new(data.Fake)
			coercedValue := model.Value{Object: map[string]interface{}{}}
			coerceToFileErr := errors.New("dummyError")

			fakeData.CoerceToFileReturns(&coercedValue, coerceToFileErr)

			objectUnderTest := _evalToFile{
				data: fakeData,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.EvalToFile(
				map[string]*model.Value{},
				map[string]interface{}{},
				new(pkg.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(coerceToFileErr))
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
			fakeInterpolater.InterpolateReturns(nil, errors.New("dummyError"))

			objectUnderTest := _evalToFile{
				interpolater: fakeInterpolater,
			}

			/* act */
			objectUnderTest.EvalToFile(
				providedScope,
				providedExpression,
				providedPkgRef,
				"dummyScratchDir",
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
				fakeInterpolater.InterpolateReturns(nil, interpolateErr)

				objectUnderTest := _evalToFile{
					interpolater: fakeInterpolater,
				}

				/* act */
				_, actualErr := objectUnderTest.EvalToFile(
					map[string]*model.Value{},
					"dummyExpression",
					new(pkg.FakeHandle),
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualErr).To(Equal(interpolateErr))

			})
		})
		Context("interpolater.Interpolate doesn't err", func() {
			It("should call data.CoerceToFile w/ expected args & return result", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				fakeInterpolater := new(interpolater.Fake)

				interpolatedValue := model.Value{String: new(string)}
				fakeInterpolater.InterpolateReturns(&interpolatedValue, nil)

				fakeData := new(data.Fake)

				coercedValue := &model.Value{Number: new(float64)}
				fakeData.CoerceToFileReturns(coercedValue, nil)

				objectUnderTest := _evalToFile{
					data:         fakeData,
					interpolater: fakeInterpolater,
				}

				/* act */
				actualFile, actualErr := objectUnderTest.EvalToFile(
					map[string]*model.Value{},
					"dummyExpression",
					new(pkg.FakeHandle),
					providedScratchDir,
				)

				/* assert */
				actualValue,
					actualScratchDir := fakeData.CoerceToFileArgsForCall(0)

				Expect(*actualValue).To(Equal(interpolatedValue))
				Expect(actualScratchDir).To(Equal(actualScratchDir))

				Expect(actualFile).To(Equal(coercedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("expression isnt float64, map[string]interface{}, or string", func() {
		It("should return expected result", func() {
			/* arrange */
			providedExpression := struct{}{}
			objectUnderTest := _evalToFile{}

			/* act */
			actualValue, actualErr := objectUnderTest.EvalToFile(
				map[string]*model.Value{},
				providedExpression,
				new(pkg.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to evaluate %+v to file; unsupported type", providedExpression)))
		})
	})
})
