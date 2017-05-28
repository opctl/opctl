package inputs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("inputs", func() {
	Context("Interpret", func() {
		Context("input param w/out arg", func() {
			Context("param is string", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "inputDefault"

						providedInputParams := map[string]*model.Param{
							"inputName": {String: &model.StringParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							validator: new(fakeValidator),
						}

						expectedInputs := map[string]*model.Data{
							providedInputName: {String: &providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							map[string]*model.Data{},
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call validator.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "inputDefault"

						providedInputParams := map[string]*model.Param{
							"inputName": {String: &model.StringParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]
						expectedInput := &model.Data{String: &providedInputDefault}

						fakeValidator := new(fakeValidator)

						objectUnderTest := _Inputs{
							validator: fakeValidator,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							map[string]*model.Data{},
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
							"inputName": {Number: &model.NumberParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							validator: new(fakeValidator),
						}

						expectedInputs := map[string]*model.Data{
							providedInputName: {Number: &providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]string{},
							providedInputParams,
							map[string]*model.Data{},
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
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

				providedScope := map[string]*model.Data{
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
					providedScope,
				)

				/* assert */
				actualArgName,
					actualArgValue,
					actualParam,
					actualScope := fakeArgInterpreter.InterpretArgsForCall(0)

				Expect(actualArgName).To(Equal(providedArgName))
				Expect(actualArgValue).To(Equal(providedArgValue))
				Expect(actualParam).To(Equal(expectedParam))
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

					expectedInput := &model.Data{String: new(string)}
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
						map[string]*model.Data{},
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
