package string

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
	"github.com/pkg/errors"
)

var _ = Context("Interpret", func() {
	It("should call interpolater.Interpolate w/ expected args & return result", func() {
		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := "dummyExpression"

		expectedErr := errors.New("dummyError")
		expectedResult := "dummyResult"
		fakeInterpolater := new(interpolater.Fake)
		fakeInterpolater.InterpolateReturns(expectedResult, expectedErr)

		expectedDeReferencer := new(interpolater.FakeDeReferencer)

		fakeDeReferencerFactory := new(fakeDeReferencerFactory)
		fakeDeReferencerFactory.NewReturns(expectedDeReferencer)

		objectUnderTest := _interpreter{
			deReferencerFactory: fakeDeReferencerFactory,
			interpolater:        fakeInterpolater,
		}

		/* act */
		actualResult, actualErr := objectUnderTest.Interpret(
			providedScope,
			providedExpression,
		)

		/* assert */
		Expect(actualResult).To(Equal(expectedResult))
		Expect(actualErr).To(Equal(expectedErr))

		actualExpression,
			actualDeReferencer := fakeInterpolater.InterpolateArgsForCall(0)

		Expect(actualExpression).To(Equal(providedExpression))
		Expect(actualDeReferencer).To(Equal(expectedDeReferencer))

	})
})
