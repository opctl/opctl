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

var _ = Context("EvalToFile", func() {
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
