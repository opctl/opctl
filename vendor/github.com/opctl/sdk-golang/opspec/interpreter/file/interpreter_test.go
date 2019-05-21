package file

import (
	"errors"
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
	"github.com/opctl/sdk-golang/opspec/interpreter/value"
)

var _ = Context("Interpret", func() {
	Context("expression is ref", func() {
		Context("expression is pkg fs ref", func() {
			It("should call interpolater.Interpolate w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}

				pkgFSRef := "/dummyPkgFSRef"
				providedExpression := fmt.Sprintf("$(%v)", pkgFSRef)
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
					"dummyScratchDir",
				)

				/* assert */
				actualExpression,
					actualScope,
					actualOpRef := fakeInterpolater.InterpolateArgsForCall(0)

				Expect(actualExpression).To(Equal(pkgFSRef))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpRef).To(Equal(providedOpRef))
			})
			Context("interpolater.Interpolate errs", func() {
				It("should return expected result", func() {

					/* arrange */
					pkgFSRef := "/dummyPkgFSRef"
					providedExpression := fmt.Sprintf("$(%v)", pkgFSRef)

					fakeInterpolater := new(interpolater.Fake)
					interpolateErr := errors.New("dummyError")
					fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

					expectedErr := fmt.Errorf(
						"unable to evaluate pkg fs ref %v; error was %v",
						pkgFSRef,
						interpolateErr.Error(),
					)

					objectUnderTest := _interpreter{
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						providedExpression,
						new(data.FakeHandle),
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(expectedErr))

				})
			})
			Context("interpolater.Interpolate doesn't error", func() {
				It("should return expected result", func() {
					/* arrange */
					providedOpHandle := new(data.FakeHandle)

					opPath := "dummyOpPath"
					providedOpHandle.PathReturns(&opPath)

					interpolatedPath := "dummyInterpolatedPath"
					fakeInterpolater := new(interpolater.Fake)
					fakeInterpolater.InterpolateReturns(interpolatedPath, nil)

					objectUnderTest := _interpreter{
						interpolater: fakeInterpolater,
					}

					expectedFileValue := filepath.Join(opPath, interpolatedPath)

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						"$(/dummyPkgFSRef)",
						providedOpHandle,
						"dummyScratchDir",
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{File: &expectedFileValue}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("expression is scope ref", func() {
			It("should call coerce.ToFile w/ expected args", func() {
				/* arrange */
				scopeName := "dummyScopeName"
				providedExpression := fmt.Sprintf("$(%v)", scopeName)
				scopeValue := model.Value{File: new(string)}
				providedScope := map[string]*model.Value{
					scopeName: &scopeValue,
				}
				providedScratchDir := "dummyScratchDir"

				fakeCoerce := new(coerce.Fake)

				objectUnderTest := _interpreter{
					coerce: fakeCoerce,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedExpression,
					new(data.FakeHandle),
					providedScratchDir,
				)

				/* assert */
				actualValue, actualScratchDir := fakeCoerce.ToFileArgsForCall(0)
				Expect(*actualValue).To(Equal(scopeValue))
				Expect(actualScratchDir).To(Equal(providedScratchDir))
			})
		})
		Context("expression is dir scope ref w/ path", func() {
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
					"dummyScratchDir",
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
						"dummyScratchDir",
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
					dirValue := "dummyDirValue"

					interpolatedPath := "dummyInterpolatedPath"
					fakeInterpolater := new(interpolater.Fake)
					fakeInterpolater.InterpolateReturns(interpolatedPath, nil)

					objectUnderTest := _interpreter{
						interpolater: fakeInterpolater,
					}

					expectedValue := filepath.Join(dirValue, interpolatedPath)

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{scopeName: {Dir: &dirValue}},
						fmt.Sprintf("$(%v/path)", scopeName),
						new(data.FakeHandle),
						"dummyScratchDir",
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{File: &expectedValue}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
	It("should call valueInterpreter.Interpret w/ expected args", func() {

		/* arrange */
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedExpression := map[string]interface{}{
			"prop1Name": "prop1Value",
		}
		providedOpRef := new(data.FakeHandle)

		fakeValueInterpreter := new(value.FakeInterpreter)
		// err to trigger immediate return
		interpretErr := errors.New("interpretErr")
		fakeValueInterpreter.InterpretReturns(model.Value{}, interpretErr)

		objectUnderTest := _interpreter{
			valueInterpreter: fakeValueInterpreter,
		}

		/* act */
		objectUnderTest.Interpret(
			providedScope,
			providedExpression,
			providedOpRef,
			"dummyScratchDir",
		)

		/* assert */
		actualExpression,
			actualScope,
			actualOpRef := fakeValueInterpreter.InterpretArgsForCall(0)

		Expect(actualExpression).To(Equal(providedExpression))
		Expect(actualScope).To(Equal(providedScope))
		Expect(actualOpRef).To(Equal(providedOpRef))

	})
	Context("valueInterpreter.Interpret errs", func() {
		It("should return expected result", func() {

			/* arrange */
			providedExpression := map[string]interface{}{
				"prop1Name": "prop1Value",
			}

			fakeValueInterpreter := new(value.FakeInterpreter)
			// err to trigger immediate return
			interpretErr := errors.New("interpretErr")
			fakeValueInterpreter.InterpretReturns(model.Value{}, interpretErr)

			expectedErr := fmt.Errorf("unable to interpret %+v to file; error was %v", providedExpression, interpretErr)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
			}

			/* act */
			_, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedExpression,
				new(data.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Context("valueInterpreter.Interpret doesn't err", func() {
		It("should call coerce.ToFile w/ expected args", func() {
			/* arrange */
			providedScratchDir := "dummyScratchDir"

			fakeValueInterpreter := new(value.FakeInterpreter)
			expectedObjectValue := model.Value{String: new(string)}
			fakeValueInterpreter.InterpretReturns(expectedObjectValue, nil)

			fakeCoerce := new(coerce.Fake)

			objectUnderTest := _interpreter{
				valueInterpreter: fakeValueInterpreter,
				coerce:           fakeCoerce,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				map[string]interface{}{},
				new(data.FakeHandle),
				providedScratchDir,
			)

			/* assert */
			actualValue,
				actualScratchDir := fakeCoerce.ToFileArgsForCall(0)
			Expect(*actualValue).To(Equal(expectedObjectValue))
			Expect(actualScratchDir).To(Equal(providedScratchDir))
		})
		It("should return expected result", func() {
			/* arrange */
			fakeCoerce := new(coerce.Fake)
			coercedValue := model.Value{Object: map[string]interface{}{}}
			toFileErr := errors.New("dummyError")

			fakeCoerce.ToFileReturns(&coercedValue, toFileErr)

			objectUnderTest := _interpreter{
				valueInterpreter: new(value.FakeInterpreter),
				coerce:           fakeCoerce,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				map[string]interface{}{},
				new(data.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(toFileErr))
		})
	})
})
