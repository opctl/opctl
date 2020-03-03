package inputs

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	inputFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/inputs/input/fakes"
	paramsFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/fakes"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("input arg", func() {
			It("should call inputFakes.Interpret w/ expected args", func() {
				/* arrange */
				providedArgName := "argName"
				providedArgValue := "argValue"

				providedInputArgs := map[string]interface{}{
					providedArgName: providedArgValue,
				}

				providedInputParams := map[string]*model.Param{
					providedArgName: {String: &model.StringParam{}},
				}

				expectedParam := providedInputParams[providedArgName]

				providedParentOpHandle := new(modelFakes.FakeDataHandle)

				providedScope := map[string]*model.Value{
					"scopeRef1Name": {},
				}

				providedOpScratchDir := "dummyOpScratchDir"

				fakeInputInterpreter := new(inputFakes.FakeInterpreter)
				// err to trigger immediate return
				fakeInputInterpreter.InterpretReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _interpreter{
					inputInterpreter: fakeInputInterpreter,
					paramsDefaulter:  new(paramsFakes.FakeDefaulter),
					paramsValidator:  new(paramsFakes.FakeValidator),
				}

				/* act */
				objectUnderTest.Interpret(
					providedInputArgs,
					providedInputParams,
					providedParentOpHandle,
					"dummyOpPath",
					providedScope,
					providedOpScratchDir,
				)

				/* assert */
				actualArgName,
					actualArgValue,
					actualParam,
					actualParentOpHandle,
					actualScope,
					actualOpScratchDir := fakeInputInterpreter.InterpretArgsForCall(0)

				Expect(actualArgName).To(Equal(providedArgName))
				Expect(actualArgValue).To(Equal(providedArgValue))
				Expect(actualParam).To(Equal(expectedParam))
				Expect(actualParentOpHandle).To(Equal(providedParentOpHandle))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpScratchDir).To(Equal(providedOpScratchDir))

			})
			Context("inputFakes.Interpret doesn't error", func() {
				It("should call paramsFakes.Default w/ expected args", func() {
					/* arrange */
					providedArgName := "argName"

					providedInputArgs := map[string]interface{}{
						providedArgName: "",
					}

					providedParams := map[string]*model.Param{
						providedArgName: nil,
					}

					providedOpPath := "dummyOpPath"

					expectedInput := &model.Value{String: new(string)}
					expectedInputs := map[string]*model.Value{
						providedArgName: expectedInput,
					}

					fakeInputInterpreter := new(inputFakes.FakeInterpreter)
					fakeInputInterpreter.InterpretReturns(expectedInput, nil)

					fakeParamsDefaulter := new(paramsFakes.FakeDefaulter)

					objectUnderTest := _interpreter{
						inputInterpreter: fakeInputInterpreter,
						paramsDefaulter:  fakeParamsDefaulter,
						paramsValidator:  new(paramsFakes.FakeValidator),
					}

					/* act */
					objectUnderTest.Interpret(
						providedInputArgs,
						providedParams,
						new(modelFakes.FakeDataHandle),
						providedOpPath,
						map[string]*model.Value{},
						"dummyOpScratchDir",
					)

					/* assert */
					actualInputs,
						actualParams,
						actualOpPath := fakeParamsDefaulter.DefaultArgsForCall(0)

					Expect(actualInputs).To(Equal(expectedInputs))
					Expect(actualParams).To(Equal(providedParams))
					Expect(actualOpPath).To(Equal(providedOpPath))
				})
				It("should call validate.Validate w/ expected args", func() {
					/* arrange */
					providedArgName := "argName"

					providedInputArgs := map[string]interface{}{
						providedArgName: "",
					}

					providedInputParams := map[string]*model.Param{
						providedArgName: nil,
					}

					fakeParamsDefaulter := new(paramsFakes.FakeDefaulter)
					expectedInputs := map[string]*model.Value{
						providedArgName: new(model.Value),
					}
					fakeParamsDefaulter.DefaultReturns(expectedInputs)

					fakeParamsValidator := new(paramsFakes.FakeValidator)

					objectUnderTest := _interpreter{
						inputInterpreter: new(inputFakes.FakeInterpreter),
						paramsDefaulter:  fakeParamsDefaulter,
						paramsValidator:  fakeParamsValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedInputArgs,
						providedInputParams,
						new(modelFakes.FakeDataHandle),
						"dummyOpPath",
						map[string]*model.Value{},
						"dummyOpScratchDir",
					)

					/* assert */
					actualInputs,
						actualParams := fakeParamsValidator.ValidateArgsForCall(0)

					Expect(actualInputs).To(Equal(expectedInputs))
					Expect(actualParams).To(Equal(providedInputParams))
				})
			})
		})
	})
})
