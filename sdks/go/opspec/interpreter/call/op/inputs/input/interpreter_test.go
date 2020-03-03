package input

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	arrayFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/array/fakes"
	booleanFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/boolean/fakes"
	dirFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/dir/fakes"
	fileFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/file/fakes"
	numberFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/number/fakes"
	objectFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/object/fakes"
	referenceFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/fakes"
	strFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/str/fakes"
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
					new(modelFakes.FakeDataHandle),
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
						new(modelFakes.FakeDataHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
		})
		Context("Arg is string", func() {
			Context("Input is array", func() {
				It("should call arrayInterpreter.Interpret w/ expected args", func() {
					/* arrange */
					providedScope := map[string]*model.Value{"dummyValue": new(model.Value)}
					providedExpression := "[dummyItem]"
					providedOpHandle := new(modelFakes.FakeDataHandle)

					fakeArrayInterpreter := new(arrayFakes.FakeInterpreter)

					fakeArrayInterpreter.InterpretReturns(nil, errors.New("dummyError"))

					objectUnderTest := _interpreter{
						arrayInterpreter: fakeArrayInterpreter,
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
						actualOpHandle := fakeArrayInterpreter.InterpretArgsForCall(0)

					Expect(actualScope).To(Equal(providedScope))
					Expect(actualExpression).To(Equal(providedExpression))
					Expect(actualOpHandle).To(Equal(providedOpHandle))

				})
				It("should return expected results", func() {
					fakeArrayInterpreter := new(arrayFakes.FakeInterpreter)

					arrayValue := new(model.Value)
					fakeArrayInterpreter.InterpretReturns(arrayValue, nil)

					objectUnderTest := _interpreter{
						arrayInterpreter: fakeArrayInterpreter,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{Array: &model.ArrayParam{}},
						new(modelFakes.FakeDataHandle),
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
					fakeBooleanInterpreter := new(booleanFakes.FakeInterpreter)

					interpolatedValue := &model.Value{Boolean: new(bool)}
					fakeBooleanInterpreter.InterpretReturns(interpolatedValue, nil)

					objectUnderTest := _interpreter{
						booleanInterpreter: fakeBooleanInterpreter,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyExpression",
						&model.Param{Boolean: &model.BooleanParam{}},
						new(modelFakes.FakeDataHandle),
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
					fakeDirInterpreter := new(dirFakes.FakeInterpreter)

					interpolatedValue := &model.Value{Dir: new(string)}
					fakeDirInterpreter.InterpretReturns(interpolatedValue, nil)

					objectUnderTest := _interpreter{
						dirInterpreter: fakeDirInterpreter,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"$(/somePkgDir)",
						&model.Param{Dir: &model.DirParam{}},
						new(modelFakes.FakeDataHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualResult).To(Equal(interpolatedValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is file", func() {
				It("should call fileInterpreter.Interpret w/ expected args", func() {
					/* arrange */
					providedScope := map[string]*model.Value{"dummyValue": new(model.Value)}
					providedExpression := "$(/somePkgFile)"
					providedOpHandle := new(modelFakes.FakeDataHandle)
					providedScratchDir := "dummyScratchDir"

					fakeFileInterpreter := new(fileFakes.FakeInterpreter)

					fakeFileInterpreter.InterpretReturns(nil, errors.New("dummyError"))

					objectUnderTest := _interpreter{
						fileInterpreter: fakeFileInterpreter,
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
						actualScratchDir := fakeFileInterpreter.InterpretArgsForCall(0)

					Expect(actualScope).To(Equal(providedScope))
					Expect(actualExpression).To(Equal(providedExpression))
					Expect(actualOpHandle).To(Equal(providedOpHandle))
					Expect(actualScratchDir).To(Equal(providedScratchDir))

				})
				It("should return expected results", func() {
					fakeFileInterpreter := new(fileFakes.FakeInterpreter)

					fileValue := new(model.Value)
					fakeFileInterpreter.InterpretReturns(fileValue, nil)

					objectUnderTest := _interpreter{
						fileInterpreter: fakeFileInterpreter,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{File: &model.FileParam{}},
						new(modelFakes.FakeDataHandle),
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
					fakeNumberInterpreter := new(numberFakes.FakeInterpreter)
					interpretedValue := model.Value{Number: new(float64)}
					fakeNumberInterpreter.InterpretReturns(&interpretedValue, nil)

					objectUnderTest := _interpreter{
						numberInterpreter: fakeNumberInterpreter,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						&model.Param{Number: &model.NumberParam{}},
						new(modelFakes.FakeDataHandle),
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

					fakeObjectInterpreter := new(objectFakes.FakeInterpreter)
					interpretedValue := model.Value{Object: new(map[string]interface{})}
					fakeObjectInterpreter.InterpretReturns(&interpretedValue, nil)

					objectUnderTest := _interpreter{
						objectInterpreter: fakeObjectInterpreter,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						providedParam,
						new(modelFakes.FakeDataHandle),
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

					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					interpretedValue := model.Value{String: new(string)}
					fakeStrInterpreter.InterpretReturns(&interpretedValue, nil)

					objectUnderTest := _interpreter{
						stringInterpreter: fakeStrInterpreter,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						providedParam,
						new(modelFakes.FakeDataHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(*actualResult).To(Equal(interpretedValue))
					Expect(actualError).To(BeNil())
				})
			})
			Context("Input is socket", func() {
				It("should return expected result", func() {
					/* arrange */
					providedParam := &model.Param{Socket: &model.SocketParam{}}

					fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
					interpretedValue := model.Value{Socket: new(string)}
					fakeReferenceInterpreter.InterpretReturns(&interpretedValue, nil)

					objectUnderTest := _interpreter{
						referenceInterpreter: fakeReferenceInterpreter,
					}

					/* act */
					actualResult, actualError := objectUnderTest.Interpret(
						"dummyName",
						"dummyValue",
						providedParam,
						new(modelFakes.FakeDataHandle),
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(*actualResult).To(Equal(interpretedValue))
					Expect(actualError).To(BeNil())
				})
				Context("referenceInterpreter.Interpret errs", func() {
					It("should return expected error", func() {
						/* arrange */
						providedName := "dummyName"

						interpolatedValue := "dummyValue"

						fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
						interpretedValue := model.Value{Socket: new(string)}
						fakeReferenceInterpreter.InterpretReturns(&interpretedValue, fmt.Errorf(""))

						objectUnderTest := _interpreter{
							referenceInterpreter: fakeReferenceInterpreter,
						}

						expectedError := fmt.Errorf("unable to bind '%v' to '%+v'; error was: ''", providedName, interpolatedValue)

						/* act */
						_, actualError := objectUnderTest.Interpret(
							providedName,
							"dummyValue",
							&model.Param{Socket: &model.SocketParam{}},
							new(modelFakes.FakeDataHandle),
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
})
