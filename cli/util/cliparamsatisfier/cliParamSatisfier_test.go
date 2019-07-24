package cliparamsatisfier

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/util/cliexiter"
	"github.com/opctl/opctl/cli/util/clioutput"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("parameterSatisfier", func() {
	Context("Satisfy", func() {
		It("should call inputSourcer.Source w/ expected args for each input", func() {
			/* arrange */
			providedInputSourcer := new(FakeInputSourcer)
			providedInputs := map[string]*types.Param{
				"input1": {String: &types.StringParam{}},
				"input2": {String: &types.StringParam{}},
			}

			expectedInputNames := map[string]struct{}{
				"input1": {},
				"input2": {},
			}

			objectUnderTest := _CLIParamSatisfier{
				cliExiter:       new(cliexiter.Fake),
				cliOutput:       new(clioutput.Fake),
				paramsValidator: new(params.FakeValidator),
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
					providedInputs := map[string]*types.Param{
						input1Name: {Array: &types.ArrayParam{}},
					}

					expectedValues := map[string]*types.Value{
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

					fakeParamsValidator := new(params.FakeValidator)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter:       new(cliexiter.Fake),
						cliOutput:       new(clioutput.Fake),
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

					providedInputs := map[string]*types.Param{
						"dummyInputName": {Boolean: &types.BooleanParam{}},
					}

					valueString := "dummyString"
					providedInputSourcer.SourceReturns(&valueString, true)

					expectedValue := types.Value{String: &valueString}

					fakeCoerce := new(coerce.Fake)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter:       new(cliexiter.Fake),
						cliOutput:       new(clioutput.Fake),
						coerce:          fakeCoerce,
						paramsValidator: new(params.FakeValidator),
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
						providedInputs := map[string]*types.Param{
							input1Name: {Boolean: &types.BooleanParam{}},
						}

						providedInputSourcer.SourceReturns(new(string), true)

						fakeCoerce := new(coerce.Fake)

						expectedBoolean := true
						expectedValues := map[string]*types.Value{
							input1Name: {
								Boolean: &expectedBoolean,
							},
						}

						fakeCoerce.ToBooleanReturns(expectedValues[input1Name], nil)

						fakeParamsValidator := new(params.FakeValidator)

						objectUnderTest := _CLIParamSatisfier{
							cliExiter:       new(cliexiter.Fake),
							cliOutput:       new(clioutput.Fake),
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

					providedInputs := map[string]*types.Param{
						"dummyInputName": {Number: &types.NumberParam{}},
					}

					valueString := "dummyString"
					providedInputSourcer.SourceReturns(&valueString, true)

					expectedValue := types.Value{String: &valueString}

					fakeCoerce := new(coerce.Fake)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter:       new(cliexiter.Fake),
						cliOutput:       new(clioutput.Fake),
						coerce:          fakeCoerce,
						paramsValidator: new(params.FakeValidator),
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
						providedInputs := map[string]*types.Param{
							input1Name: {Number: &types.NumberParam{}},
						}

						providedInputSourcer.SourceReturns(new(string), true)

						fakeCoerce := new(coerce.Fake)

						expectedNumber := 2.2
						expectedValues := map[string]*types.Value{
							input1Name: {
								Number: &expectedNumber,
							},
						}

						fakeCoerce.ToNumberReturns(expectedValues[input1Name], nil)

						fakeParamsValidator := new(params.FakeValidator)

						objectUnderTest := _CLIParamSatisfier{
							cliExiter:       new(cliexiter.Fake),
							cliOutput:       new(clioutput.Fake),
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
					providedInputs := map[string]*types.Param{
						input1Name: {Object: &types.ObjectParam{}},
					}

					expectedValues := map[string]*types.Value{
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

					fakeParamsValidator := new(params.FakeValidator)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter:       new(cliexiter.Fake),
						cliOutput:       new(clioutput.Fake),
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
