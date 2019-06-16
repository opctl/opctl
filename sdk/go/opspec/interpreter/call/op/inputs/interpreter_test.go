package inputs

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/data"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/op/inputs/input"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/op/params"
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
			It("should call input.Interpret w/ expected args", func() {
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

				providedParentOpHandle := new(data.FakeHandle)

				providedScope := map[string]*model.Value{
					"scopeRef1Name": {},
				}

				providedOpScratchDir := "dummyOpScratchDir"

				fakeInputInterpreter := new(input.FakeInterpreter)
				// err to trigger immediate return
				fakeInputInterpreter.InterpretReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _interpreter{
					inputInterpreter: fakeInputInterpreter,
					paramsDefaulter:  new(params.FakeDefaulter),
					paramsValidator:  new(params.FakeValidator),
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
			Context("input.Interpret doesn't error", func() {
				It("should call params.Default w/ expected args", func() {
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

					fakeInputInterpreter := new(input.FakeInterpreter)
					fakeInputInterpreter.InterpretReturns(expectedInput, nil)

					fakeParamsDefaulter := new(params.FakeDefaulter)

					objectUnderTest := _interpreter{
						inputInterpreter: fakeInputInterpreter,
						paramsDefaulter:  fakeParamsDefaulter,
						paramsValidator:  new(params.FakeValidator),
					}

					/* act */
					objectUnderTest.Interpret(
						providedInputArgs,
						providedParams,
						new(data.FakeHandle),
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

					fakeParamsDefaulter := new(params.FakeDefaulter)
					expectedInputs := map[string]*model.Value{
						providedArgName: new(model.Value),
					}
					fakeParamsDefaulter.DefaultReturns(expectedInputs)

					fakeParamsValidator := new(params.FakeValidator)

					objectUnderTest := _interpreter{
						inputInterpreter: new(input.FakeInterpreter),
						paramsDefaulter:  fakeParamsDefaulter,
						paramsValidator:  fakeParamsValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedInputArgs,
						providedInputParams,
						new(data.FakeHandle),
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
