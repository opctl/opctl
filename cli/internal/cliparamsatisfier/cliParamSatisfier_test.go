package cliparamsatisfier

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	clioutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	. "github.com/opctl/opctl/cli/internal/cliparamsatisfier/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/data/coerce/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	paramsFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/fakes"
)

var _ = Context("parameterSatisfier", func() {
	Context("Satisfy", func() {
		It("should call inputSourcer.Source w/ expected args for each input", func() {
			/* arrange */
			providedInputSourcer := new(FakeInputSourcer)
			providedInputs := map[string]*model.Param{
				"input1": {String: &model.StringParam{}},
				"input2": {String: &model.StringParam{}},
			}

			expectedInputNames := map[string]struct{}{
				"input1": {},
				"input2": {},
			}

			objectUnderTest := _CLIParamSatisfier{
				cliExiter:       new(cliexiterFakes.FakeCliExiter),
				cliOutput:       new(clioutputFakes.FakeCliOutput),
				paramsValidator: new(paramsFakes.FakeValidator),
			}

			/* act */
			objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

			/* assert */
			actualInputNames := map[string]struct{}{}
			for callIndex := 0; callIndex < providedInputSourcer.SourceCallCount(); callIndex++ {
				actualInputName := providedInputSourcer.SourceArgsForCall(callIndex)
				actualInputNames[actualInputName] = struct{}{}
			}

			Expect(actualInputNames).To(Equal(expectedInputNames))
		})
		Context("param.Array isn't nil", func() {
			Context("value isn't nil", func() {

				It("should call inputs.validate w/ expected args", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)

					input1Name := "input1Name"
					providedInputs := map[string]*model.Param{
						input1Name: {Array: &model.ArrayParam{}},
					}

					expectedValues := map[string]*model.Value{
						input1Name: {
							Array: new([]interface{}),
						},
					}

					valueBytes, err := json.Marshal(expectedValues[input1Name].Array)
					if nil != err {
						Fail(err.Error())
					}

					valueString := string(valueBytes)
					providedInputSourcer.SourceReturns(&valueString, true)

					fakeParamsValidator := new(paramsFakes.FakeValidator)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter:       new(cliexiterFakes.FakeCliExiter),
						cliOutput:       new(clioutputFakes.FakeCliOutput),
						paramsValidator: fakeParamsValidator,
					}

					/* act */
					objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					actualValues,
						actualParams := fakeParamsValidator.ValidateArgsForCall(0)

					Expect(actualValues).To(Equal(expectedValues))
					Expect(actualParams).To(Equal(providedInputs))
				})
			})
		})
		Context("param.Boolean isn't nil", func() {
			Context("value isn't nil", func() {
				It("should call data.CoerceToBoolean w/ expected args", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)

					providedInputs := map[string]*model.Param{
						"dummyInputName": {Boolean: &model.BooleanParam{}},
					}

					valueString := "dummyString"
					providedInputSourcer.SourceReturns(&valueString, true)

					expectedValue := model.Value{String: &valueString}

					fakeCoerce := new(FakeCoerce)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter:       new(cliexiterFakes.FakeCliExiter),
						cliOutput:       new(clioutputFakes.FakeCliOutput),
						coerce:          fakeCoerce,
						paramsValidator: new(paramsFakes.FakeValidator),
					}

					/* act */
					objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					actualValue := fakeCoerce.ToBooleanArgsForCall(0)
					Expect(*actualValue).To(Equal(expectedValue))
				})
				Context("data.CoerceToBoolean doesn't err", func() {
					It("should call inputs.validate w/ expected args", func() {
						/* arrange */
						providedInputSourcer := new(FakeInputSourcer)

						input1Name := "input1Name"
						providedInputs := map[string]*model.Param{
							input1Name: {Boolean: &model.BooleanParam{}},
						}

						providedInputSourcer.SourceReturns(new(string), true)

						fakeCoerce := new(FakeCoerce)

						expectedBoolean := true
						expectedValues := map[string]*model.Value{
							input1Name: {
								Boolean: &expectedBoolean,
							},
						}

						fakeCoerce.ToBooleanReturns(expectedValues[input1Name], nil)

						fakeParamsValidator := new(paramsFakes.FakeValidator)

						objectUnderTest := _CLIParamSatisfier{
							cliExiter:       new(cliexiterFakes.FakeCliExiter),
							cliOutput:       new(clioutputFakes.FakeCliOutput),
							coerce:          fakeCoerce,
							paramsValidator: fakeParamsValidator,
						}

						/* act */
						objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

						/* assert */
						actualValues, actualParams := fakeParamsValidator.ValidateArgsForCall(0)
						Expect(actualValues).To(Equal(expectedValues))
						Expect(actualParams).To(Equal(providedInputs))
					})
				})
			})
		})
		Context("param.Number isn't nil", func() {
			Context("value isn't nil", func() {
				It("should call data.CoerceToNumber w/ expected args", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)

					providedInputs := map[string]*model.Param{
						"dummyInputName": {Number: &model.NumberParam{}},
					}

					valueString := "dummyString"
					providedInputSourcer.SourceReturns(&valueString, true)

					expectedValue := model.Value{String: &valueString}

					fakeCoerce := new(FakeCoerce)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter:       new(cliexiterFakes.FakeCliExiter),
						cliOutput:       new(clioutputFakes.FakeCliOutput),
						coerce:          fakeCoerce,
						paramsValidator: new(paramsFakes.FakeValidator),
					}

					/* act */
					objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					actualValue := fakeCoerce.ToNumberArgsForCall(0)
					Expect(*actualValue).To(Equal(expectedValue))
				})
				Context("data.CoerceToNumber doesn't err", func() {
					It("should call inputs.validate w/ expected args", func() {
						/* arrange */
						providedInputSourcer := new(FakeInputSourcer)

						input1Name := "input1Name"
						providedInputs := map[string]*model.Param{
							input1Name: {Number: &model.NumberParam{}},
						}

						providedInputSourcer.SourceReturns(new(string), true)

						fakeCoerce := new(FakeCoerce)

						expectedNumber := 2.2
						expectedValues := map[string]*model.Value{
							input1Name: {
								Number: &expectedNumber,
							},
						}

						fakeCoerce.ToNumberReturns(expectedValues[input1Name], nil)

						fakeParamsValidator := new(paramsFakes.FakeValidator)

						objectUnderTest := _CLIParamSatisfier{
							cliExiter:       new(cliexiterFakes.FakeCliExiter),
							cliOutput:       new(clioutputFakes.FakeCliOutput),
							coerce:          fakeCoerce,
							paramsValidator: fakeParamsValidator,
						}

						/* act */
						objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

						/* assert */
						actualValues, actualParams := fakeParamsValidator.ValidateArgsForCall(0)
						Expect(actualValues).To(Equal(expectedValues))
						Expect(actualParams).To(Equal(providedInputs))
					})
				})
			})
		})
		Context("param.Object isn't nil", func() {
			Context("value isn't nil", func() {

				It("should call inputs.validate w/ expected args", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)

					input1Name := "input1Name"
					providedInputs := map[string]*model.Param{
						input1Name: {Object: &model.ObjectParam{}},
					}

					expectedValues := map[string]*model.Value{
						input1Name: {
							Object: new(map[string]interface{}),
						},
					}

					valueBytes, err := json.Marshal(expectedValues[input1Name].Object)
					if nil != err {
						Fail(err.Error())
					}

					valueString := string(valueBytes)
					providedInputSourcer.SourceReturns(&valueString, true)

					fakeParamsValidator := new(paramsFakes.FakeValidator)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter:       new(cliexiterFakes.FakeCliExiter),
						cliOutput:       new(clioutputFakes.FakeCliOutput),
						paramsValidator: fakeParamsValidator,
					}

					/* act */
					objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					actualValues,
						actualParams := fakeParamsValidator.ValidateArgsForCall(0)

					Expect(actualValues).To(Equal(expectedValues))
					Expect(actualParams).To(Equal(providedInputs))
				})
			})
		})
	})
})
