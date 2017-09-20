package expression

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/pkg/errors"
	"path/filepath"
)

var _ = Context("EvalToDir", func() {
	Context("expression is scope ref", func() {
		It("should return expected result", func() {
			/* arrange */
			scopeName := "dummyScopeName"
			providedExpression := fmt.Sprintf("$(%v)", scopeName)
			scopeValue := model.Value{Dir: new(string)}
			providedScope := map[string]*model.Value{
				scopeName: &scopeValue,
			}

			objectUnderTest := _dirEvaluator{}

			/* act */
			actualValue, actualErr := objectUnderTest.EvalToDir(
				providedScope,
				providedExpression,
				new(pkg.FakeHandle),
			)

			/* assert */
			Expect(*actualValue).To(Equal(scopeValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("expression is deprecated pkg fs ref", func() {
		Context("interpolater.Interpolate errors", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeInterpolater := new(interpolater.Fake)
				interpolateError := errors.New("dummyError")
				fakeInterpolater.InterpolateReturns("", interpolateError)

				objectUnderTest := _dirEvaluator{
					interpolater: fakeInterpolater,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.EvalToDir(
					map[string]*model.Value{},
					"/deprecatedPkgFsRef",
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(interpolateError))
			})
		})
		Context("interpolater.Interpolate doesn't error", func() {
			It("should return expected result", func() {
				/* arrange */
				interpolatedExpression := "dummyExpression"

				fakeInterpolater := new(interpolater.Fake)
				fakeInterpolater.InterpolateReturns(interpolatedExpression, nil)

				objectUnderTest := _dirEvaluator{
					interpolater: fakeInterpolater,
				}

				fakeHandle := new(pkg.FakeHandle)
				pkgRef := "dummyPkgRef"
				fakeHandle.RefReturns(pkgRef)

				expectedFileValue := filepath.Join(pkgRef, interpolatedExpression)

				/* act */
				actualValue, actualErr := objectUnderTest.EvalToDir(
					map[string]*model.Value{},
					"/deprecatedPkgFsRef",
					fakeHandle,
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Dir: &expectedFileValue}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	It("should call interpolater.Interpolate w/ expected args", func() {
		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := "dummyExpression"
		providedPkgRef := new(pkg.FakeHandle)

		fakeInterpolater := new(interpolater.Fake)
		// err to trigger immediate return
		fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

		objectUnderTest := _dirEvaluator{
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

		interpolateErr := errors.New("dummyError")
		fakeInterpolater.InterpolateReturns("", interpolateErr)

		objectUnderTest := _dirEvaluator{
			interpolater: fakeInterpolater,
		}

		/* act */
		actualValue, actualErr := objectUnderTest.EvalToDir(
			map[string]*model.Value{},
			"dummyExpression",
			new(pkg.FakeHandle),
		)

		/* assert */
		Expect(actualValue).To(BeNil())
		Expect(actualErr).To(Equal(interpolateErr))

	})
})
