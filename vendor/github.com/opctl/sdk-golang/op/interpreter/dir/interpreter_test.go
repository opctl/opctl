package dir

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/op/interpreter/interpolater"
	"path/filepath"
)

var _ = Context("Interpret", func() {
	Context("expression is scope ref", func() {
		It("should return expected result", func() {
			/* arrange */
			scopeName := "dummyScopeName"
			providedExpression := fmt.Sprintf("$(%v)", scopeName)
			scopeValue := model.Value{Dir: new(string)}
			providedScope := map[string]*model.Value{
				scopeName: &scopeValue,
			}

			objectUnderTest := _interpreter{}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				providedScope,
				providedExpression,
				new(data.FakeHandle),
			)

			/* assert */
			Expect(*actualValue).To(Equal(scopeValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("expression is scope ref w/ path", func() {
		It("should call interpolater.Interpolate w/ expected args", func() {
			/* arrange */
			scopeName := "dummyScopeName"
			providedScope := map[string]*model.Value{scopeName: {Dir: new(string)}}

			providedPath := "dummyPath"
			providedExpression := fmt.Sprintf("$(%v/%v)", scopeName, providedPath)
			providedOpRef := new(data.FakeHandle)

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
				providedOpRef,
			)

			/* assert */
			actualExpression,
				actualScope,
				actualOpRef := fakeInterpolater.InterpolateArgsForCall(0)

			Expect(actualExpression).To(Equal(providedPath))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpRef).To(Equal(providedOpRef))
		})
		Context("interpolater.Interpolate errs", func() {
			It("should return expected result", func() {

				/* arrange */
				scopeName := "dummyScopeName"
				providedScope := map[string]*model.Value{scopeName: {Dir: new(string)}}

				providedPath := "dummyPath"
				providedExpression := fmt.Sprintf("$(%v/%v)", scopeName, providedPath)

				fakeInterpolater := new(interpolater.Fake)
				interpolateErr := errors.New("dummyError")
				fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

				expectedErr := fmt.Errorf(
					"unable to evaluate path %v; error was %v",
					providedPath,
					interpolateErr.Error(),
				)

				objectUnderTest := _interpreter{
					interpolater: fakeInterpolater,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					providedScope,
					providedExpression,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(expectedErr))

			})
		})
		Context("interpolater.Interpolate doesn't error", func() {
			It("should return expected result", func() {
				/* arrange */
				scopeName := "dummyScopeName"
				scopeValue := "dummyScopeValue"

				interpolatedExpression := "dummyInterpolatedExpression"
				fakeInterpolater := new(interpolater.Fake)
				fakeInterpolater.InterpolateReturns(interpolatedExpression, nil)

				objectUnderTest := _interpreter{
					interpolater: fakeInterpolater,
				}

				expectedFileValue := filepath.Join(scopeValue, interpolatedExpression)

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{scopeName: {Dir: &scopeValue}},
					fmt.Sprintf("$(%v/path)", scopeName),
					new(data.FakeHandle),
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Dir: &expectedFileValue}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
