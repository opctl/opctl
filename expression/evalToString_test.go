package expression

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/interpolater"
	"github.com/pkg/errors"
)

var _ = Context("EvalToString", func() {
	It("should call interpolater.Interpolate w/ expected args & return result", func() {
		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := "dummyExpression"
		providedPkgHandle := new(pkg.FakeHandle)

		expectedErr := errors.New("dummyError")
		expectedResult := "dummyResult"
		fakeInterpolater := new(interpolater.Fake)
		fakeInterpolater.InterpolateReturns(expectedResult, expectedErr)

		objectUnderTest := _evalToString{
			interpolater: fakeInterpolater,
		}

		/* act */
		actualResult, actualErr := objectUnderTest.EvalToString(
			providedScope,
			providedExpression,
			providedPkgHandle,
		)

		/* assert */
		Expect(actualResult).To(Equal(expectedResult))
		Expect(actualErr).To(Equal(expectedErr))

		actualExpression,
			actualScope,
			actualPkgHandle := fakeInterpolater.InterpolateArgsForCall(0)

		Expect(actualExpression).To(Equal(providedExpression))
		Expect(actualScope).To(Equal(providedScope))
		Expect(actualPkgHandle).To(Equal(providedPkgHandle))

	})
})
