package inputs

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/params"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("inputs.interpreter", func() {
	Context("input arg", func() {
		It("should call argInterpreter.Interpret w/ expected args", func() {
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

			providedParentPkgHandle := new(pkg.FakeHandle)

			providedScope := map[string]*model.Value{
				"scopeRef1Name": {},
			}

			providedOpScratchDir := "dummyOpScratchDir"

			fakeArgInterpreter := new(fakeArgInterpreter)
			// err to trigger immediate return
			fakeArgInterpreter.InterpretReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _interpreter{
				argInterpreter: fakeArgInterpreter,
				params:         new(params.Fake),
				data:           new(data.Fake),
			}

			/* act */
			objectUnderTest.Interpret(
				providedInputArgs,
				providedInputParams,
				providedParentPkgHandle,
				"dummyPkgPath",
				providedScope,
				providedOpScratchDir,
			)

			/* assert */
			actualArgName,
				actualArgValue,
				actualParam,
				actualParentPkgHandle,
				actualScope,
				actualOpScratchDir := fakeArgInterpreter.InterpretArgsForCall(0)

			Expect(actualArgName).To(Equal(providedArgName))
			Expect(actualArgValue).To(Equal(providedArgValue))
			Expect(actualParam).To(Equal(expectedParam))
			Expect(actualParentPkgHandle).To(Equal(providedParentPkgHandle))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpScratchDir).To(Equal(providedOpScratchDir))

		})
		Context("argInterpreter.Interpret doesn't error", func() {
			It("should call params.Default w/ expected args", func() {
				/* arrange */
				providedArgName := "argName"

				providedInputArgs := map[string]interface{}{
					providedArgName: "",
				}

				providedParams := map[string]*model.Param{
					providedArgName: nil,
				}

				providedPkgPath := "dummyPkgPath"

				expectedInput := &model.Value{String: new(string)}
				expectedInputs := map[string]*model.Value{
					providedArgName: expectedInput,
				}

				fakeArgInterpreter := new(fakeArgInterpreter)
				fakeArgInterpreter.InterpretReturns(expectedInput, nil)

				fakeParams := new(params.Fake)

				objectUnderTest := _interpreter{
					argInterpreter: fakeArgInterpreter,
					params:         fakeParams,
					data:           new(data.Fake),
				}

				/* act */
				objectUnderTest.Interpret(
					providedInputArgs,
					providedParams,
					new(pkg.FakeHandle),
					providedPkgPath,
					map[string]*model.Value{},
					"dummyOpScratchDir",
				)

				/* assert */
				actualInputs,
					actualParams,
					actualPkgPath := fakeParams.DefaultArgsForCall(0)

				Expect(actualInputs).To(Equal(expectedInputs))
				Expect(actualParams).To(Equal(providedParams))
				Expect(actualPkgPath).To(Equal(providedPkgPath))
			})
			It("should call data.Validate w/ expected args", func() {
				/* arrange */
				providedArgName := "argName"

				providedInputArgs := map[string]interface{}{
					providedArgName: "",
				}

				providedInputParams := map[string]*model.Param{
					providedArgName: nil,
				}

				expectedParam := providedInputParams[providedArgName]

				expectedInput := &model.Value{String: new(string)}
				fakeParams := new(params.Fake)
				fakeParams.DefaultReturns(map[string]*model.Value{
					providedArgName: expectedInput,
				})

				fakeData := new(data.Fake)

				objectUnderTest := _interpreter{
					argInterpreter: new(fakeArgInterpreter),
					params:         fakeParams,
					data:           fakeData,
				}

				/* act */
				objectUnderTest.Interpret(
					providedInputArgs,
					providedInputParams,
					new(pkg.FakeHandle),
					"dummyPkgPath",
					map[string]*model.Value{},
					"dummyOpScratchDir",
				)

				/* assert */
				actualInput,
					actualParam := fakeData.ValidateArgsForCall(0)

				Expect(actualInput).To(Equal(expectedInput))
				Expect(actualParam).To(Equal(expectedParam))
			})
		})
	})
})
