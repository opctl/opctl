package expression

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/interpolater"
	"github.com/pkg/errors"
)

var _ = Context("EvalToDir", func() {
	It("should call interpolater.Interpolate w/ expected args", func() {
		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := "dummyExpression"
		providedPkgRef := new(pkg.FakeHandle)

		fakeInterpolater := new(interpolater.Fake)
		// err to trigger immediate return
		fakeInterpolater.InterpolateReturns(nil, errors.New("dummyError"))

		objectUnderTest := _evalToDir{
			interpolater: fakeInterpolater,
		}

		/* act */
		objectUnderTest.EvalToDir(
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
	It("should return expected result", func() {
		/* arrange */
		fakeInterpolater := new(interpolater.Fake)

		interpolateValue := &model.Value{Dir: new(string)}
		interpolateErr := errors.New("dummyError")
		fakeInterpolater.InterpolateReturns(interpolateValue, interpolateErr)

		objectUnderTest := _evalToDir{
			interpolater: fakeInterpolater,
		}

		/* act */
		actualValue, actualErr := objectUnderTest.EvalToDir(
			map[string]*model.Value{},
			"dummyExpression",
			new(pkg.FakeHandle),
		)

		/* assert */
		Expect(actualValue).To(Equal(interpolateValue))
		Expect(actualErr).To(Equal(interpolateErr))

	})
})
