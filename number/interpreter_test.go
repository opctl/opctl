package number

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
	"github.com/pkg/errors"
)

var _ = Context("Interpret", func() {
	It("should call interpolater.Interpolate w/ expected args", func() {
		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := "dummyExpression"

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
		)

		/* assert */
		actualExpression,
			actualScope := fakeInterpolater.InterpolateArgsForCall(0)

		Expect(actualExpression).To(Equal(providedExpression))
		Expect(actualScope).To(Equal(providedScope))

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
			)

			/* assert */
			Expect(actualErr).To(Equal(interpolateErr))

		})
	})
	Context("interpolater.Interpolate doesn't err", func() {
		It("should call data.CoerceToNumber w/ expected args & return result", func() {
			/* arrange */
			fakeInterpolater := new(interpolater.Fake)

			interpolatedData := "dummyInterpolatedData"
			fakeInterpolater.InterpolateReturns(interpolatedData, nil)

			fakeData := new(data.Fake)

			coercedValue := 2.2
			fakeData.CoerceToNumberReturns(coercedValue, nil)

			objectUnderTest := _interpreter{
				data:         fakeData,
				interpolater: fakeInterpolater,
			}

			/* act */
			actualNumber, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				"dummyExpression",
			)

			/* assert */
			Expect(*fakeData.CoerceToNumberArgsForCall(0)).To(Equal(model.Value{String: &interpolatedData}))

			Expect(actualNumber).To(Equal(coercedValue))
			Expect(actualErr).To(BeNil())
		})
	})
})
