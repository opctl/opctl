package inputs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

var _ = Context("inputs", func() {
	Context("Interpret", func() {
		Context("input param w/out arg", func() {
			Context("param is dir", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "/pkgDirDefault"
						providedPkgPath := "dummyPkgPath"

						providedInputParams := map[string]*model.Param{
							providedInputName: {Dir: &model.DirParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							validator: new(fakeValidator),
						}

						expectedInputValue := filepath.Join(providedPkgPath, providedInputDefault)
						expectedInputs := map[string]*model.Value{
							providedInputName: {Dir: &expectedInputValue},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							providedPkgPath,
							map[string]*model.Value{},
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call validator.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "/pkgDirDefault"
						providedPkgPath := "dummyPkgPath"

						providedInputParams := map[string]*model.Param{
							providedInputName: {Dir: &model.DirParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]

						expectedInputValue := filepath.Join(providedPkgPath, providedInputDefault)
						expectedInput := &model.Value{Dir: &expectedInputValue}

						fakeValidator := new(fakeValidator)

						objectUnderTest := _Inputs{
							validator: fakeValidator,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							providedPkgPath,
							map[string]*model.Value{},
						)

						/* assert */
						actualInput, actualParam := fakeValidator.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is file", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "/pkgFileDefault"
						providedPkgPath := "dummyPkgPath"

						providedInputParams := map[string]*model.Param{
							providedInputName: {File: &model.FileParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							validator: new(fakeValidator),
						}

						expectedInputValue := filepath.Join(providedPkgPath, providedInputDefault)
						expectedInputs := map[string]*model.Value{
							providedInputName: {File: &expectedInputValue},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							providedPkgPath,
							map[string]*model.Value{},
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call validator.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "/pkgFileDefault"
						providedPkgPath := "pkgPath"

						providedInputParams := map[string]*model.Param{
							providedInputName: {File: &model.FileParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]

						expectedInputValue := filepath.Join(providedPkgPath, providedInputDefault)
						expectedInput := &model.Value{File: &expectedInputValue}

						fakeValidator := new(fakeValidator)

						objectUnderTest := _Inputs{
							validator: fakeValidator,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							providedPkgPath,
							map[string]*model.Value{},
						)

						/* assert */
						actualInput, actualParam := fakeValidator.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is number", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := 2.2

						providedInputParams := map[string]*model.Param{
							providedInputName: {Number: &model.NumberParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							validator: new(fakeValidator),
						}

						expectedInputs := map[string]*model.Value{
							providedInputName: {Number: &providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							"dummyPkgPath",
							map[string]*model.Value{},
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call validator.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := 2.2

						providedInputParams := map[string]*model.Param{
							providedInputName: {Number: &model.NumberParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]
						expectedInput := &model.Value{Number: &providedInputDefault}

						fakeValidator := new(fakeValidator)

						objectUnderTest := _Inputs{
							validator: fakeValidator,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							"dummyPkgPath",
							map[string]*model.Value{},
						)

						/* assert */
						actualInput, actualParam := fakeValidator.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is object", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := map[string]interface{}{}

						providedInputParams := map[string]*model.Param{
							providedInputName: {Object: &model.ObjectParam{Default: providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							validator: new(fakeValidator),
						}

						expectedInputs := map[string]*model.Value{
							providedInputName: {Object: providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							"dummyPkgPath",
							map[string]*model.Value{},
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call validator.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := map[string]interface{}{}

						providedInputParams := map[string]*model.Param{
							providedInputName: {Object: &model.ObjectParam{Default: providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]
						expectedInput := &model.Value{Object: providedInputDefault}

						fakeValidator := new(fakeValidator)

						objectUnderTest := _Inputs{
							validator: fakeValidator,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							"dummyPkgPath",
							map[string]*model.Value{},
						)

						/* assert */
						actualInput, actualParam := fakeValidator.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is string", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "inputDefault"

						providedInputParams := map[string]*model.Param{
							providedInputName: {String: &model.StringParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							validator: new(fakeValidator),
						}

						expectedInputs := map[string]*model.Value{
							providedInputName: {String: &providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							"dummyPkgPath",
							map[string]*model.Value{},
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call validator.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "inputDefault"

						providedInputParams := map[string]*model.Param{
							providedInputName: {String: &model.StringParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]
						expectedInput := &model.Value{String: &providedInputDefault}

						fakeValidator := new(fakeValidator)

						objectUnderTest := _Inputs{
							validator: fakeValidator,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							"dummyPkgPath",
							map[string]*model.Value{},
						)

						/* assert */
						actualInput, actualParam := fakeValidator.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
		})
		Context("input arg", func() {
			It("should call argInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedArgName := "argName"
				providedArgValue := "argValue"

				providedInputArgs := map[string]string{
					providedArgName: providedArgValue,
				}

				providedInputParams := map[string]*model.Param{
					providedArgName: {String: &model.StringParam{}},
				}

				expectedParam := providedInputParams[providedArgName]

				providedPkgPath := "pkgPath"

				providedScope := map[string]*model.Value{
					"scopeRef1Name": {},
				}

				fakeArgInterpreter := new(fakeArgInterpreter)

				objectUnderTest := _Inputs{
					argInterpreter: fakeArgInterpreter,
					validator:      new(fakeValidator),
				}

				/* act */
				objectUnderTest.Interpret(
					providedInputArgs,
					providedInputParams,
					providedPkgPath,
					providedScope,
				)

				/* assert */
				actualArgName,
					actualArgValue,
					actualParam,
					actualPkgPath,
					actualScope := fakeArgInterpreter.InterpretArgsForCall(0)

				Expect(actualArgName).To(Equal(providedArgName))
				Expect(actualArgValue).To(Equal(providedArgValue))
				Expect(actualParam).To(Equal(expectedParam))
				Expect(actualPkgPath).To(Equal(providedPkgPath))
				Expect(actualScope).To(Equal(providedScope))
			})
			Context("argInterpreter.Interpret doesn't error", func() {
				It("should call validator.Validate w/ expected args", func() {
					/* arrange */
					providedArgName := "argName"

					providedInputArgs := map[string]string{
						providedArgName: "",
					}

					providedInputParams := map[string]*model.Param{
						providedArgName: nil,
					}

					expectedParam := providedInputParams[providedArgName]

					expectedInput := &model.Value{String: new(string)}
					fakeArgInterpreter := new(fakeArgInterpreter)
					fakeArgInterpreter.InterpretReturns(expectedInput, nil)

					fakeValidator := new(fakeValidator)

					objectUnderTest := _Inputs{
						argInterpreter: fakeArgInterpreter,
						validator:      fakeValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedInputArgs,
						providedInputParams,
						"dummyPkgPath",
						map[string]*model.Value{},
					)

					/* assert */
					actualInput, actualParam := fakeValidator.ValidateArgsForCall(0)

					Expect(actualInput).To(Equal(expectedInput))
					Expect(actualParam).To(Equal(expectedParam))
				})
			})
		})
	})
})
