package value

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("Interpret", func() {
	var _ = Context("Interpret", func() {
		Context("expression is bool", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValueExpression := false
				objectUnderTest := _interpreter{}

				/* act */
				actualValue, _ := objectUnderTest.Interpret(
					providedValueExpression,
					map[string]*types.Value{},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(types.Value{Boolean: &providedValueExpression}))
			})
		})
		Context("expression is float64", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValueExpression := 2.2
				objectUnderTest := _interpreter{}

				/* act */
				actualValue, _ := objectUnderTest.Interpret(
					providedValueExpression,
					map[string]*types.Value{},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(types.Value{Number: &providedValueExpression}))
			})
		})
		Context("expression is int", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValueExpression := 2
				expectedNumber := float64(providedValueExpression)
				objectUnderTest := _interpreter{}

				/* act */
				actualValue, _ := objectUnderTest.Interpret(
					providedValueExpression,
					map[string]*types.Value{},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(types.Value{Number: &expectedNumber}))
			})
		})
		Context("expression is map[string]interface{}", func() {
		})
		Context("expression is []interface{}", func() {
		})
		Context("expression is string", func() {
			It("should call interpolater.Interpolate w/ expected args & return expected result", func() {
				/* arrange */
				providedScope := map[string]*types.Value{"dummyName": {}}
				providedExpression := "dummyExpression"
				providedOpRef := new(data.FakeHandle)

				fakeInterpolater := new(interpolater.Fake)
				// err to trigger immediate return
				fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

				objectUnderTest := _interpreter{
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.Interpret(
					providedExpression,
					providedScope,
					providedOpRef,
				)

				/* assert */
				actualExpression,
					actualScope,
					actualOpRef := fakeInterpolater.InterpolateArgsForCall(0)

				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpRef).To(Equal(providedOpRef))

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
						"dummyExpression",
						map[string]*types.Value{},
						new(data.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(interpolateErr))

				})
			})
		})
		It("should return expected result", func() {
			/* arrange */
			fakeInterpolater := new(interpolater.Fake)

			interpolatedValue := "dummyString"
			fakeInterpolater.InterpolateReturns(interpolatedValue, nil)

			objectUnderTest := _interpreter{
				interpolater: fakeInterpolater,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				"dummyExpression",
				map[string]*types.Value{},
				new(data.FakeHandle),
			)

			/* assert */
			Expect(actualValue).To(Equal(types.Value{String: &interpolatedValue}))
			Expect(actualErr).To(BeNil())
		})
	})
})
