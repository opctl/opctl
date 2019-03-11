package file

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
	arrayInitializer "github.com/opctl/sdk-golang/opspec/interpreter/array/initializer"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
	objectInitializer "github.com/opctl/sdk-golang/opspec/interpreter/object/initializer"
	"path/filepath"
)

var _ = Context("Interpret", func() {
	Context("expression is float64", func() {
		It("should call coerce.ToFile w/ expected args", func() {
			/* arrange */
			providedExpression := 2.2
			providedScratchDir := "dummyScratchDir"

			fakeCoerce := new(coerce.Fake)

			objectUnderTest := _interpreter{
				coerce: fakeCoerce,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedExpression,
				new(data.FakeHandle),
				providedScratchDir,
			)

			/* assert */
			actualValue,
				actualScratchDir := fakeCoerce.ToFileArgsForCall(0)
			Expect(*actualValue).To(Equal(model.Value{Number: &providedExpression}))
			Expect(actualScratchDir).To(Equal(providedScratchDir))
		})
		It("should return expected result", func() {
			/* arrange */
			fakeCoerce := new(coerce.Fake)
			coercedValue := model.Value{Number: new(float64)}
			toFileErr := errors.New("dummyError")

			fakeCoerce.ToFileReturns(&coercedValue, toFileErr)

			objectUnderTest := _interpreter{
				coerce: fakeCoerce,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				2.2,
				new(data.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(toFileErr))
		})
	})
	Context("expression is int", func() {
		It("should call coerce.ToFile w/ expected args", func() {
			/* arrange */
			providedExpression := 2
			providedScratchDir := "dummyScratchDir"

			expectedNumber := float64(providedExpression)

			fakeCoerce := new(coerce.Fake)

			objectUnderTest := _interpreter{
				coerce: fakeCoerce,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedExpression,
				new(data.FakeHandle),
				providedScratchDir,
			)

			/* assert */
			actualValue,
				actualScratchDir := fakeCoerce.ToFileArgsForCall(0)
			Expect(*actualValue).To(Equal(model.Value{Number: &expectedNumber}))
			Expect(actualScratchDir).To(Equal(providedScratchDir))
		})
		It("should return expected result", func() {
			/* arrange */
			fakeCoerce := new(coerce.Fake)
			coercedValue := model.Value{Number: new(float64)}
			toFileErr := errors.New("dummyError")

			fakeCoerce.ToFileReturns(&coercedValue, toFileErr)

			objectUnderTest := _interpreter{
				coerce: fakeCoerce,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				2,
				new(data.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(*actualValue).To(Equal(coercedValue))
			Expect(actualErr).To(Equal(toFileErr))
		})
	})
	Context("expression is map[string]interface{}", func() {
		It("should call objectInitializerInterpreter.Interpret w/ expected args", func() {

			/* arrange */
			providedScope := map[string]*model.Value{"dummyName": {}}
			providedExpression := map[string]interface{}{
				"prop1Name": "prop1Value",
			}
			providedOpRef := new(data.FakeHandle)

			fakeObjectInitializerInterpreter := new(objectInitializer.FakeInterpreter)
			// err to trigger immediate return
			interpretErr := errors.New("interpretErr")
			fakeObjectInitializerInterpreter.InterpretReturns(map[string]interface{}{}, interpretErr)

			objectUnderTest := _interpreter{
				objectInitializerInterpreter: fakeObjectInitializerInterpreter,
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
				actualOpRef := fakeObjectInitializerInterpreter.InterpretArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpRef).To(Equal(providedOpRef))

		})
		Context("objectInitializerInterpreter.Interpret errs", func() {
			It("should return expected result", func() {

				/* arrange */
				providedExpression := map[string]interface{}{
					"prop1Name": "prop1Value",
				}

				fakeObjectInitializerInterpreter := new(objectInitializer.FakeInterpreter)
				// err to trigger immediate return
				interpretErr := errors.New("interpretErr")
				fakeObjectInitializerInterpreter.InterpretReturns(map[string]interface{}{}, interpretErr)

				expectedErr := fmt.Errorf("unable to evaluate %+v to file; error was %v", providedExpression, interpretErr)

				objectUnderTest := _interpreter{
					objectInitializerInterpreter: fakeObjectInitializerInterpreter,
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
		Context("objectInitializerInterpreter.Interpret doesn't err", func() {
			It("should call coerce.ToFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				fakeObjectInitializerInterpreter := new(objectInitializer.FakeInterpreter)
				expectedObjectValue := map[string]interface{}{"dummyName": 2.2}
				fakeObjectInitializerInterpreter.InterpretReturns(expectedObjectValue, nil)

				fakeCoerce := new(coerce.Fake)

				objectUnderTest := _interpreter{
					objectInitializerInterpreter: fakeObjectInitializerInterpreter,
					coerce: fakeCoerce,
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
				Expect(*actualValue).To(Equal(model.Value{Object: expectedObjectValue}))
				Expect(actualScratchDir).To(Equal(providedScratchDir))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeCoerce := new(coerce.Fake)
				coercedValue := model.Value{Object: map[string]interface{}{}}
				toFileErr := errors.New("dummyError")

				fakeCoerce.ToFileReturns(&coercedValue, toFileErr)

				objectUnderTest := _interpreter{
					objectInitializerInterpreter: new(objectInitializer.FakeInterpreter),
					coerce: fakeCoerce,
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
	Context("expression is []interface{}", func() {
		It("should call arrayInitializerInterpreter.Interpret w/ expected args", func() {

			/* arrange */
			providedScope := map[string]*model.Value{"dummyName": {}}
			providedExpression := []interface{}{
				"item1",
			}
			providedOpRef := new(data.FakeHandle)

			fakeArrayInitializerInterpreter := new(arrayInitializer.FakeInterpreter)
			// err to trigger immediate return
			interpretErr := errors.New("interpretErr")
			fakeArrayInitializerInterpreter.InterpretReturns([]interface{}{}, interpretErr)

			objectUnderTest := _interpreter{
				arrayInitializerInterpreter: fakeArrayInitializerInterpreter,
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
				actualOpRef := fakeArrayInitializerInterpreter.InterpretArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpRef).To(Equal(providedOpRef))

		})
		Context("arrayInitializerInterpreter.Interpret errs", func() {
			It("should return expected result", func() {

				/* arrange */
				providedExpression := []interface{}{
					"item1",
				}

				fakeArrayInitializerInterpreter := new(arrayInitializer.FakeInterpreter)
				// err to trigger immediate return
				interpretErr := errors.New("interpretErr")
				fakeArrayInitializerInterpreter.InterpretReturns([]interface{}{}, interpretErr)

				expectedErr := fmt.Errorf("unable to evaluate %+v to file; error was %v", providedExpression, interpretErr)

				objectUnderTest := _interpreter{
					arrayInitializerInterpreter: fakeArrayInitializerInterpreter,
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
		Context("arrayInitializerInterpreter.Interpret doesn't err", func() {
			It("should call coerce.ToFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				fakeArrayInitializerInterpreter := new(arrayInitializer.FakeInterpreter)
				expectedArrayValue := []interface{}{"item1"}
				fakeArrayInitializerInterpreter.InterpretReturns(expectedArrayValue, nil)

				fakeCoerce := new(coerce.Fake)

				objectUnderTest := _interpreter{
					arrayInitializerInterpreter: fakeArrayInitializerInterpreter,
					coerce: fakeCoerce,
				}

				/* act */
				objectUnderTest.Interpret(
					map[string]*model.Value{},
					[]interface{}{},
					new(data.FakeHandle),
					providedScratchDir,
				)

				/* assert */
				actualValue,
					actualScratchDir := fakeCoerce.ToFileArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Array: expectedArrayValue}))
				Expect(actualScratchDir).To(Equal(providedScratchDir))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeCoerce := new(coerce.Fake)
				coercedValue := model.Value{Array: []interface{}{}}
				toFileErr := errors.New("dummyError")

				fakeCoerce.ToFileReturns(&coercedValue, toFileErr)

				objectUnderTest := _interpreter{
					arrayInitializerInterpreter: new(arrayInitializer.FakeInterpreter),
					coerce: fakeCoerce,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					[]interface{}{},
					new(data.FakeHandle),
					"dummyScratchDir",
				)

				/* assert */
				Expect(*actualValue).To(Equal(coercedValue))
				Expect(actualErr).To(Equal(toFileErr))
			})
		})
	})
	Context("expression is string", func() {
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
		Context("expression starts w/ scope ref", func() {
			It("should call interpolater.Interpolate w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{"dummyKey": {String: new(string)}}
				providedExpression := "$(dummyRef)suffix"
				providedOpHandle := new(data.FakeHandle)

				interpolatedPath := "dummyInterpolatedPath"
				fakeInterpolater := new(interpolater.Fake)
				// err to trigger immediate return
				fakeInterpolater.InterpolateReturns(interpolatedPath, errors.New("dummyError"))

				objectUnderTest := _interpreter{
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedExpression,
					providedOpHandle,
					"dummyScratchDir",
				)

				/* assert */
				actualExpression,
					actualScope,
					actualOpHandle := fakeInterpolater.InterpolateArgsForCall(0)

				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
			})
			Context("interpolater.Interpolate errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := "$(dummyRef)suffix"

					fakeInterpolater := new(interpolater.Fake)
					interpolateErr := errors.New("dummyError")
					fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

					expectedErr := fmt.Errorf(
						"unable to evaluate %v to file; error was %v",
						providedExpression,
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
			Context("interpolater.Interpolate doesn't err", func() {
				It("should call coerce.ToFile w/ expected args & return result", func() {
					/* arrange */
					providedExpression := "$(dummyRef)suffix"
					providedScratchDir := "dummyScratchDir"

					interpolatedPath := "dummyInterpolatedPath"
					fakeInterpolater := new(interpolater.Fake)
					fakeInterpolater.InterpolateReturns(interpolatedPath, nil)

					fakeCoerce := new(coerce.Fake)
					coerceValue := model.Value{File: new(string)}
					coerceErr := errors.New("dummyErr")
					fakeCoerce.ToFileReturns(&coerceValue, coerceErr)

					objectUnderTest := _interpreter{
						coerce:       fakeCoerce,
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						providedExpression,
						new(data.FakeHandle),
						providedScratchDir,
					)

					/* assert */
					actualCoerceValueArg, actualCoerceScratchDirArg := fakeCoerce.ToFileArgsForCall(0)
					Expect(*actualCoerceValueArg).To(Equal(model.Value{String: &interpolatedPath}))
					Expect(actualCoerceScratchDirArg).To(Equal(providedScratchDir))

					Expect(*actualValue).To(Equal(coerceValue))
					Expect(actualErr).To(Equal(coerceErr))
				})
			})
		})
		Context("expression isn't ref", func() {
			It("should call interpolater.Interpolate w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{"dummyKey": {String: new(string)}}
				providedExpression := "dummyExpression"
				providedOpHandle := new(data.FakeHandle)

				interpolatedPath := "dummyInterpolatedPath"
				fakeInterpolater := new(interpolater.Fake)
				// err to trigger immediate return
				fakeInterpolater.InterpolateReturns(interpolatedPath, errors.New("dummyError"))

				objectUnderTest := _interpreter{
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedExpression,
					providedOpHandle,
					"dummyScratchDir",
				)

				/* assert */
				actualExpression,
					actualScope,
					actualOpHandle := fakeInterpolater.InterpolateArgsForCall(0)

				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
			})
			Context("interpolater.Interpolate errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := "dummyExpression"

					fakeInterpolater := new(interpolater.Fake)
					interpolateErr := errors.New("dummyError")
					fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

					expectedErr := fmt.Errorf(
						"unable to evaluate %v to file; error was %v",
						providedExpression,
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
			Context("interpolater.Interpolate doesn't err", func() {
				It("should call coerce.ToFile w/ expected args & return result", func() {
					/* arrange */
					providedExpression := "dummyExpression"
					providedScratchDir := "dummyScratchDir"

					interpolatedPath := "dummyInterpolatedPath"
					fakeInterpolater := new(interpolater.Fake)
					fakeInterpolater.InterpolateReturns(interpolatedPath, nil)

					fakeCoerce := new(coerce.Fake)
					coerceValue := model.Value{File: new(string)}
					coerceErr := errors.New("dummyErr")
					fakeCoerce.ToFileReturns(&coerceValue, coerceErr)

					objectUnderTest := _interpreter{
						coerce:       fakeCoerce,
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						providedExpression,
						new(data.FakeHandle),
						providedScratchDir,
					)

					/* assert */
					actualCoerceValueArg, actualCoerceScratchDirArg := fakeCoerce.ToFileArgsForCall(0)
					Expect(*actualCoerceValueArg).To(Equal(model.Value{String: &interpolatedPath}))
					Expect(actualCoerceScratchDirArg).To(Equal(providedScratchDir))

					Expect(*actualValue).To(Equal(coerceValue))
					Expect(actualErr).To(Equal(coerceErr))
				})
			})
		})
	})
	Context("expression isnt float64, map[string]interface{}, or string", func() {
		It("should return expected result", func() {
			/* arrange */
			providedExpression := struct{}{}
			objectUnderTest := _interpreter{}

			/* act */
			actualValue, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedExpression,
				new(data.FakeHandle),
				"dummyScratchDir",
			)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to evaluate %+v to file", providedExpression)))
		})
	})
})
