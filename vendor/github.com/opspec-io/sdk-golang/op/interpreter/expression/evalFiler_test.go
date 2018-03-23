package expression

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
	"path/filepath"
)

var _ = Context("EvalToFile", func() {
	Context("expression is float64", func() {
		It("should call coerce.ToFile w/ expected args", func() {
			/* arrange */
			providedExpression := 2.2
			providedScratchDir := "dummyScratchDir"

			fakeCoerce := new(coerce.Fake)

			objectUnderTest := _evalFiler{
				coerce: fakeCoerce,
			}

			/* act */
			objectUnderTest.EvalToFile(
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

			objectUnderTest := _evalFiler{
				coerce: fakeCoerce,
			}

			/* act */
			actualValue, actualErr := objectUnderTest.EvalToFile(
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
	Context("expression is map[string]interface{}", func() {
		It("should call evalObjectInitializerer.Eval w/ expected args", func() {

			/* arrange */
			providedScope := map[string]*model.Value{"dummyName": {}}
			providedExpression := map[string]interface{}{
				"prop1Name": "prop1Value",
			}
			providedOpRef := new(data.FakeHandle)

			fakeEvalObjectInitializerer := new(fakeEvalObjectInitializerer)
			// err to trigger immediate return
			evalErr := errors.New("evalErr")
			fakeEvalObjectInitializerer.EvalReturns(map[string]interface{}{}, evalErr)

			objectUnderTest := _evalFiler{
				evalObjectInitializerer: fakeEvalObjectInitializerer,
			}

			/* act */
			objectUnderTest.EvalToFile(
				providedScope,
				providedExpression,
				providedOpRef,
				"dummyScratchDir",
			)

			/* assert */
			actualExpression,
				actualScope,
				actualOpRef := fakeEvalObjectInitializerer.EvalArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpRef).To(Equal(providedOpRef))

		})
		Context("evalObjectInitializerer.Eval errs", func() {
			It("should return expected result", func() {

				/* arrange */
				providedExpression := map[string]interface{}{
					"prop1Name": "prop1Value",
				}

				fakeEvalObjectInitializerer := new(fakeEvalObjectInitializerer)
				// err to trigger immediate return
				evalErr := errors.New("evalErr")
				fakeEvalObjectInitializerer.EvalReturns(map[string]interface{}{}, evalErr)

				expectedErr := fmt.Errorf("unable to evaluate %+v to file; error was %v", providedExpression, evalErr)

				objectUnderTest := _evalFiler{
					evalObjectInitializerer: fakeEvalObjectInitializerer,
				}

				/* act */
				_, actualErr := objectUnderTest.EvalToFile(
					map[string]*model.Value{},
					providedExpression,
					new(data.FakeHandle),
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("evalObjectInitializerer.Eval doesn't err", func() {
			It("should call coerce.ToFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				fakeEvalObjectInitializerer := new(fakeEvalObjectInitializerer)
				expectedObjectValue := map[string]interface{}{"dummyName": 2.2}
				fakeEvalObjectInitializerer.EvalReturns(expectedObjectValue, nil)

				fakeCoerce := new(coerce.Fake)

				objectUnderTest := _evalFiler{
					evalObjectInitializerer: fakeEvalObjectInitializerer,
					coerce:                  fakeCoerce,
				}

				/* act */
				objectUnderTest.EvalToFile(
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

				objectUnderTest := _evalFiler{
					evalObjectInitializerer: new(fakeEvalObjectInitializerer),
					coerce:                  fakeCoerce,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.EvalToFile(
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
		It("should call evalArrayInitializerer.Eval w/ expected args", func() {

			/* arrange */
			providedScope := map[string]*model.Value{"dummyName": {}}
			providedExpression := []interface{}{
				"item1",
			}
			providedOpRef := new(data.FakeHandle)

			fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
			// err to trigger immediate return
			evalErr := errors.New("evalErr")
			fakeEvalArrayInitializerer.EvalReturns([]interface{}{}, evalErr)

			arrayUnderTest := _evalFiler{
				evalArrayInitializerer: fakeEvalArrayInitializerer,
			}

			/* act */
			arrayUnderTest.EvalToFile(
				providedScope,
				providedExpression,
				providedOpRef,
				"dummyScratchDir",
			)

			/* assert */
			actualExpression,
				actualScope,
				actualOpRef := fakeEvalArrayInitializerer.EvalArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpRef).To(Equal(providedOpRef))

		})
		Context("evalArrayInitializerer.Eval errs", func() {
			It("should return expected result", func() {

				/* arrange */
				providedExpression := []interface{}{
					"item1",
				}

				fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
				// err to trigger immediate return
				evalErr := errors.New("evalErr")
				fakeEvalArrayInitializerer.EvalReturns([]interface{}{}, evalErr)

				expectedErr := fmt.Errorf("unable to evaluate %+v to file; error was %v", providedExpression, evalErr)

				arrayUnderTest := _evalFiler{
					evalArrayInitializerer: fakeEvalArrayInitializerer,
				}

				/* act */
				_, actualErr := arrayUnderTest.EvalToFile(
					map[string]*model.Value{},
					providedExpression,
					new(data.FakeHandle),
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("evalArrayInitializerer.Eval doesn't err", func() {
			It("should call coerce.ToFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				fakeEvalArrayInitializerer := new(fakeEvalArrayInitializerer)
				expectedArrayValue := []interface{}{"item1"}
				fakeEvalArrayInitializerer.EvalReturns(expectedArrayValue, nil)

				fakeCoerce := new(coerce.Fake)

				arrayUnderTest := _evalFiler{
					evalArrayInitializerer: fakeEvalArrayInitializerer,
					coerce:                 fakeCoerce,
				}

				/* act */
				arrayUnderTest.EvalToFile(
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

				arrayUnderTest := _evalFiler{
					evalArrayInitializerer: new(fakeEvalArrayInitializerer),
					coerce:                 fakeCoerce,
				}

				/* act */
				actualValue, actualErr := arrayUnderTest.EvalToFile(
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

				objectUnderTest := _evalFiler{
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.EvalToFile(
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

					objectUnderTest := _evalFiler{
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToFile(
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

					objectUnderTest := _evalFiler{
						interpolater: fakeInterpolater,
					}

					expectedFileValue := filepath.Join(opPath, interpolatedPath)

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToFile(
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

				objectUnderTest := _evalFiler{
					coerce: fakeCoerce,
				}

				/* act */
				objectUnderTest.EvalToFile(
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

				objectUnderTest := _evalFiler{
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.EvalToFile(
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

					objectUnderTest := _evalFiler{
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToFile(
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

					objectUnderTest := _evalFiler{
						interpolater: fakeInterpolater,
					}

					expectedValue := filepath.Join(dirValue, interpolatedPath)

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToFile(
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

				objectUnderTest := _evalFiler{
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.EvalToFile(
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

					objectUnderTest := _evalFiler{
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToFile(
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

					objectUnderTest := _evalFiler{
						coerce:       fakeCoerce,
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToFile(
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

				objectUnderTest := _evalFiler{
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.EvalToFile(
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

					objectUnderTest := _evalFiler{
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToFile(
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

					objectUnderTest := _evalFiler{
						coerce:       fakeCoerce,
						interpolater: fakeInterpolater,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.EvalToFile(
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
			objectUnderTest := _evalFiler{}

			/* act */
			actualValue, actualErr := objectUnderTest.EvalToFile(
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
