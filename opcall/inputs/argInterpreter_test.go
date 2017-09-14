package inputs

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
)

var _ = Context("argInterpreter", func() {
	Context("Interpret", func() {
		Context("param nil", func() {
			It("should return expected error", func() {
				/* arrange */
				providedName := "dummyName"

				expectedError := fmt.Errorf("unable to bind to '%v'; '%v' not a defined input", providedName, providedName)

				objectUnderTest := _argInterpreter{}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					"dummyName",
					"dummyValue",
					nil,
					new(pkg.FakeHandle),
					map[string]*model.Value{},
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("Implicit arg", func() {
			Context("Ref not in scope", func() {
				It("should return expected error", func() {
					/* arrange */
					providedName := "dummyName"

					expectedError := fmt.Errorf("unable to bind to '%v' via implicit ref; '%v' not in scope", providedName, providedName)

					objectUnderTest := _argInterpreter{}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						providedName,
						"",
						&model.Param{},
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("Ref in scope", func() {
				It("should call validate.Validate w/ expected args & return result", func() {
					/* arrange */
					providedName := "dummyName"
					providedValue := ""
					providedParam := &model.Param{}
					expectedValue := &model.Value{String: new(string)}
					providedScope := map[string]*model.Value{providedName: expectedValue}

					objectUnderTest := _argInterpreter{}

					/* act */
					actualValue, actualError := objectUnderTest.Interpret(
						providedName,
						providedValue,
						providedParam,
						new(pkg.FakeHandle),
						providedScope,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualValue).To(Equal(expectedValue))
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("Deprecated explicit arg", func() {
			It("should call validate.Validate w/ expected args & return result", func() {
				/* arrange */
				providedValue := "dummyValue"
				providedParam := &model.Param{}
				expectedValue := &model.Value{String: new(string)}
				providedScope := map[string]*model.Value{providedValue: expectedValue}

				objectUnderTest := _argInterpreter{}

				/* act */
				actualValue, actualError := objectUnderTest.Interpret(
					"dummyName",
					providedValue,
					providedParam,
					new(pkg.FakeHandle),
					providedScope,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualValue).To(Equal(expectedValue))
				Expect(actualError).To(BeNil())
			})
		})
		Context("Explicit arg", func() {
			Context("Ref not in scope", func() {
				It("should return expected error", func() {
					/* arrange */
					providedName := "dummyName"
					explicitRef := "dummyRef"
					providedValue := fmt.Sprintf("$(%v)", explicitRef)
					providedParam := &model.Param{}

					expectedError := fmt.Errorf("unable to bind '%v' to '%v'", providedName, providedValue)

					objectUnderTest := _argInterpreter{}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						providedName,
						providedValue,
						providedParam,
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("Ref in scope", func() {
				It("should call validate.Validate w/ expected args & return result", func() {
					/* arrange */
					explicitRef := "dummyRef"
					providedValue := fmt.Sprintf("$(%v)", explicitRef)
					providedParam := &model.Param{}
					expectedValue := &model.Value{String: new(string)}
					providedScope := map[string]*model.Value{explicitRef: expectedValue}

					objectUnderTest := _argInterpreter{}

					/* act */
					actualValue, actualError := objectUnderTest.Interpret(
						"dummyName",
						providedValue,
						providedParam,
						new(pkg.FakeHandle),
						providedScope,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualValue).To(Equal(expectedValue))
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("Interpolated arg", func() {
			Context("Input is string", func() {
				It("should return expected result", func() {
					/* arrange */
					providedParam := &model.Param{String: &model.StringParam{}}

					fakeExpression := new(expression.Fake)
					interpretedValue := "dummyValue"
					fakeExpression.EvalToStringReturns(interpretedValue, nil)

					expectedResult := &model.Value{String: &interpretedValue}

					objectUnderTest := _argInterpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						providedParam,
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is Dir", func() {
				Context("bound to pkg dir", func() {
					It("should return expected results", func() {
						/* arrange */
						fakeExpression := new(expression.Fake)

						interpolatedValue := "dummyValue"
						fakeExpression.EvalToStringReturns(interpolatedValue, nil)

						expectedResult := &model.Value{Dir: &interpolatedValue}

						objectUnderTest := _argInterpreter{
							expression: fakeExpression,
						}

						/* act */
						actualResult, actualError := objectUnderTest.Interpret(
							"dummyName",
							"$(/somePkgDir)",
							&model.Param{Dir: &model.DirParam{}},
							new(pkg.FakeHandle),
							map[string]*model.Value{},
							"dummyScratchDir",
						)

						/* assert */
						Expect(actualResult).To(Equal(expectedResult))
						Expect(actualError).To(BeNil())
					})
				})
				It("should return expected result", func() {
					/* arrange */
					fakeExpression := new(expression.Fake)
					interpolatedValue := "dummyValue"
					fakeExpression.EvalToStringReturns(interpolatedValue, nil)

					expectedResult := &model.Value{Dir: &interpolatedValue}

					objectUnderTest := _argInterpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{Dir: &model.DirParam{}},
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
					Expect(actualError).To(BeNil())
				})
				It("should root path", func() {
					/* arrange */
					fakeExpression := new(expression.Fake)

					expectedValue := fmt.Sprintf("%v%v", string(filepath.Separator), "dummyValue")
					interpolatedValue := fmt.Sprintf("..\\../%v../..\\", expectedValue)
					fakeExpression.EvalToStringReturns(interpolatedValue, nil)

					expectedResult := &model.Value{Dir: &expectedValue}

					objectUnderTest := _argInterpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{Dir: &model.DirParam{}},
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is Number", func() {
				It("should call validate w/ expected args", func() {
					/* arrange */
					fakeExpression := new(expression.Fake)
					interpretedValue := float64(2.1)
					fakeExpression.EvalToNumberReturns(interpretedValue, nil)

					expectedResult := &model.Value{Number: &interpretedValue}

					objectUnderTest := _argInterpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{Number: &model.NumberParam{}},
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is File", func() {
				It("should call expression.EvalToFile w/ expected args", func() {
					/* arrange */
					providedScope := map[string]*model.Value{"dummyValue": new(model.Value)}
					providedExpression := "$(/somePkgFile)"
					providedPkgHandle := new(pkg.FakeHandle)
					providedScratchDir := "dummyScratchDir"

					fakeExpression := new(expression.Fake)

					fakeExpression.EvalToFileReturns(nil, errors.New("dummyError"))

					objectUnderTest := _argInterpreter{
						expression: fakeExpression,
					}

					/* act */
					objectUnderTest.Interpret(
						"dummyName",
						providedExpression,
						&model.Param{File: &model.FileParam{}},
						providedPkgHandle,
						providedScope,
						providedScratchDir,
					)

					/* assert */
					actualScope,
						actualExpression,
						actualPkgHandle,
						actualScratchDir := fakeExpression.EvalToFileArgsForCall(0)

					Expect(actualScope).To(Equal(providedScope))
					Expect(actualExpression).To(Equal(providedExpression))
					Expect(actualPkgHandle).To(Equal(providedPkgHandle))
					Expect(actualScratchDir).To(Equal(providedScratchDir))

				})
				It("should return expected results", func() {
					fakeExpression := new(expression.Fake)

					fileValue := new(model.Value)
					fakeExpression.EvalToFileReturns(fileValue, nil)

					objectUnderTest := _argInterpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{File: &model.FileParam{}},
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(fileValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is Socket", func() {
				It("should return expected error", func() {
					/* arrange */
					providedName := "dummyName"
					fakeExpression := new(expression.Fake)

					interpolatedValue := "dummyValue"
					fakeExpression.EvalToStringReturns(interpolatedValue, nil)

					expectedError := fmt.Errorf("unable to bind '%v' to '%v'; sockets must be passed by reference", providedName, interpolatedValue)

					objectUnderTest := _argInterpreter{
						expression: fakeExpression,
					}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						&model.Param{Socket: &model.SocketParam{}},
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
		})
	})
})
