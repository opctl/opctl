package value

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater/fakes"
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
					map[string]*model.Value{},
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(model.Value{Boolean: &providedValueExpression}))
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
					map[string]*model.Value{},
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(model.Value{Number: &providedValueExpression}))
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
					map[string]*model.Value{},
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(model.Value{Number: &expectedNumber}))
			})
		})
		Context("expression is map[string]interface{}", func() {
		})
		Context("expression is []interface{}", func() {
		})
		Context("expression is string", func() {
			It("should call interpolaterFakes.Interpolate w/ expected args & return expected result", func() {
				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedExpression := "dummyExpression"
				providedOpRef := new(modelFakes.FakeDataHandle)

				fakeInterpolater := new(FakeInterpolater)
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
			Context("interpolaterFakes.Interpolate errs", func() {
				It("should return expected err", func() {
					/* arrange */
					fakeInterpolater := new(FakeInterpolater)
					interpolateErr := errors.New("dummyError")
					fakeInterpolater.InterpolateReturns("", interpolateErr)

					objectUnderTest := _interpreter{
						interpolater: fakeInterpolater,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						"dummyExpression",
						map[string]*model.Value{},
						new(modelFakes.FakeDataHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(interpolateErr))

				})
			})
		})
		It("should return expected result", func() {
			/* arrange */
			fakeInterpolater := new(FakeInterpolater)

			interpolatedValue := "dummyString"
			fakeInterpolater.InterpolateReturns(interpolatedValue, nil)

			objectUnderTest := _interpreter{
				interpolater: fakeInterpolater,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				"dummyExpression",
				map[string]*model.Value{},
				new(modelFakes.FakeDataHandle),
			)

			/* assert */
			Expect(actualValue).To(Equal(model.Value{String: &interpolatedValue}))
			Expect(actualErr).To(BeNil())
		})
	})
})
