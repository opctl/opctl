package input

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("param nil", func() {
			It("should return expected error", func() {
				/* arrange */
				providedName := "dummyName"

				expectedError := fmt.Errorf("unable to bind to '%v'; '%v' not a defined input", providedName, providedName)

				objectUnderTest := _interpreter{}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					"dummyName",
					"dummyValue",
					nil,
					new(data.FakeHandle),
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

					objectUnderTest := _interpreter{}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						providedName,
						"",
						&model.Param{},
						new(data.FakeHandle),
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

					objectUnderTest := _interpreter{}

					/* act */
					actualValue, actualError := objectUnderTest.Interpret(
						providedName,
						providedValue,
						providedParam,
						new(data.FakeHandle),
						providedScope,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualValue).To(Equal(expectedValue))
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("Arg is string", func() {
			Context("Input is array", func() {
				It("should call expression.EvalToArray w/ expected args", func() {
					/* arrange */
					providedScope := map[string]*model.Value{"dummyValue": new(model.Value)}
					providedExpression := "[dummyItem]"
					providedOpHandle := new(data.FakeHandle)

					fakeExpression := new(expression.Fake)

					fakeExpression.EvalToArrayReturns(nil, errors.New("dummyError"))

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					objectUnderTest.Interpret(
						"dummyName",
						providedExpression,
						&model.Param{Array: &model.ArrayParam{}},
						providedOpHandle,
						providedScope,
						"dummyScratchDir",
					)

					/* assert */
					actualScope,
						actualExpression,
						actualOpHandle := fakeExpression.EvalToArrayArgsForCall(0)

					Expect(actualScope).To(Equal(providedScope))
					Expect(actualExpression).To(Equal(providedExpression))
					Expect(actualOpHandle).To(Equal(providedOpHandle))

				})
				It("should return expected results", func() {
					fakeExpression := new(expression.Fake)

					arrayValue := new(model.Value)
					fakeExpression.EvalToArrayReturns(arrayValue, nil)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{Array: &model.ArrayParam{}},
						new(data.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(arrayValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is boolean", func() {
				It("should return expected results", func() {
					/* arrange */
					fakeExpression := new(expression.Fake)

					interpolatedValue := &model.Value{Boolean: new(bool)}
					fakeExpression.EvalToBooleanReturns(interpolatedValue, nil)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyExpression",
						&model.Param{Boolean: &model.BooleanParam{}},
						new(data.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(interpolatedValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is dir", func() {
				It("should return expected results", func() {
					/* arrange */
					fakeExpression := new(expression.Fake)

					interpolatedValue := &model.Value{Dir: new(string)}
					fakeExpression.EvalToDirReturns(interpolatedValue, nil)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"$(/somePkgDir)",
						&model.Param{Dir: &model.DirParam{}},
						new(data.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(interpolatedValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is file", func() {
				It("should call expression.EvalToFile w/ expected args", func() {
					/* arrange */
					providedScope := map[string]*model.Value{"dummyValue": new(model.Value)}
					providedExpression := "$(/somePkgFile)"
					providedOpHandle := new(data.FakeHandle)
					providedScratchDir := "dummyScratchDir"

					fakeExpression := new(expression.Fake)

					fakeExpression.EvalToFileReturns(nil, errors.New("dummyError"))

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					objectUnderTest.Interpret(
						"dummyName",
						providedExpression,
						&model.Param{File: &model.FileParam{}},
						providedOpHandle,
						providedScope,
						providedScratchDir,
					)

					/* assert */
					actualScope,
						actualExpression,
						actualOpHandle,
						actualScratchDir := fakeExpression.EvalToFileArgsForCall(0)

					Expect(actualScope).To(Equal(providedScope))
					Expect(actualExpression).To(Equal(providedExpression))
					Expect(actualOpHandle).To(Equal(providedOpHandle))
					Expect(actualScratchDir).To(Equal(providedScratchDir))

				})
				It("should return expected results", func() {
					fakeExpression := new(expression.Fake)

					fileValue := new(model.Value)
					fakeExpression.EvalToFileReturns(fileValue, nil)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{File: &model.FileParam{}},
						new(data.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(fileValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is number", func() {
				It("should call validate w/ expected args", func() {
					/* arrange */
					fakeExpression := new(expression.Fake)
					interpretedValue := model.Value{Number: new(float64)}
					fakeExpression.EvalToNumberReturns(&interpretedValue, nil)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{Number: &model.NumberParam{}},
						new(data.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(*actualResult).To(Equal(interpretedValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is object", func() {
				It("should return expected result", func() {
					/* arrange */
					providedParam := &model.Param{Object: &model.ObjectParam{}}

					fakeExpression := new(expression.Fake)
					interpretedValue := model.Value{Object: map[string]interface{}{"dummyName": "dummyValue"}}
					fakeExpression.EvalToObjectReturns(&interpretedValue, nil)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						providedParam,
						new(data.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(*actualResult).To(Equal(interpretedValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is string", func() {
				It("should return expected result", func() {
					/* arrange */
					providedParam := &model.Param{String: &model.StringParam{}}

					fakeExpression := new(expression.Fake)
					interpretedValue := model.Value{String: new(string)}
					fakeExpression.EvalToStringReturns(&interpretedValue, nil)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						providedParam,
						new(data.FakeHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(*actualResult).To(Equal(interpretedValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is socket", func() {
				It("should return expected error", func() {
					/* arrange */
					providedName := "dummyName"
					fakeExpression := new(expression.Fake)

					interpolatedValue := "dummyValue"
					interpretedValue := model.Value{String: &interpolatedValue}
					fakeExpression.EvalToStringReturns(&interpretedValue, nil)

					expectedError := fmt.Errorf("unable to bind '%v' to '%+v'; sockets must be passed by reference", providedName, interpolatedValue)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						&model.Param{Socket: &model.SocketParam{}},
						new(data.FakeHandle),
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
