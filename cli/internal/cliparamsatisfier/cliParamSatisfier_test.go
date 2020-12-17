package cliparamsatisfier

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clioutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	. "github.com/opctl/opctl/cli/internal/cliparamsatisfier/internal/fakes"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("parameterSatisfier", func() {
	Context("Satisfy", func() {
		It("should call inputSourcer.Source w/ expected args for each input", func() {
			/* arrange */
			providedInputSourcer := new(FakeInputSourcer)
			providedInputSourcer.SourceReturns(nil, true)
			providedInputs := map[string]*model.Param{
				"input1": {String: &model.StringParam{}},
				"input2": {String: &model.StringParam{}},
			}

			expectedInputNames := map[string]struct{}{
				"input1": {},
				"input2": {},
			}

			objectUnderTest := _CLIParamSatisfier{
				cliOutput: new(clioutputFakes.FakeCliOutput),
			}

			/* act */
			_, err := objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

			/* assert */
			actualInputNames := map[string]struct{}{}
			for callIndex := 0; callIndex < providedInputSourcer.SourceCallCount(); callIndex++ {
				actualInputName := providedInputSourcer.SourceArgsForCall(callIndex)
				actualInputNames[actualInputName] = struct{}{}
			}

			Expect(err).To(BeNil())
			Expect(actualInputNames).To(Equal(expectedInputNames))
		})
		Context("param.Array isn't nil", func() {
			Context("value isn't nil", func() {

				It("should return expected outputs", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)

					input1Name := "input1Name"
					providedInputs := map[string]*model.Param{
						input1Name: {Array: &model.ArrayParam{}},
					}

					expectedOutputs := map[string]*model.Value{
						input1Name: {
							Array: new([]interface{}),
						},
					}

					valueBytes, err := json.Marshal(expectedOutputs[input1Name].Array)
					if nil != err {
						Fail(err.Error())
					}

					valueString := string(valueBytes)
					providedInputSourcer.SourceReturns(&valueString, true)

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: new(clioutputFakes.FakeCliOutput),
					}

					/* act */
					actualOutputs, err := objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					Expect(err).To(BeNil())
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param.Boolean isn't nil", func() {
			Context("value isn't nil", func() {
				It("should return expected outputs", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)
					inputIdentifier := "inputIdentifier"

					providedInputs := map[string]*model.Param{
						inputIdentifier: {Boolean: &model.BooleanParam{}},
					}

					valueBool := true
					valueString := fmt.Sprintf("%v", valueBool)
					providedInputSourcer.SourceReturns(&valueString, true)

					expectedOutputs := map[string]*model.Value{
						inputIdentifier: &model.Value{Boolean: &valueBool},
					}

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: new(clioutputFakes.FakeCliOutput),
					}

					/* act */
					actualOutputs, err := objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					Expect(err).To(BeNil())
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param.Number isn't nil", func() {
			Context("value isn't nil", func() {
				It("should call data.CoerceToNumber w/ expected args", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)
					inputIdentifier := "inputIdentifier"

					providedInputs := map[string]*model.Param{
						inputIdentifier: {Number: &model.NumberParam{}},
					}

					valueNumber := 1.1
					valueString := fmt.Sprintf("%v", valueNumber)
					providedInputSourcer.SourceReturns(&valueString, true)

					expectedOutputs := map[string]*model.Value{
						inputIdentifier: &model.Value{Number: &valueNumber},
					}

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: new(clioutputFakes.FakeCliOutput),
					}

					/* act */
					actualOutputs, err := objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					Expect(err).To(BeNil())
					Expect(actualOutputs).To(Equal(expectedOutputs))
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

					expectedOutputs := map[string]*model.Value{
						input1Name: {
							Object: new(map[string]interface{}),
						},
					}

					valueBytes, err := json.Marshal(expectedOutputs[input1Name].Object)
					if nil != err {
						Fail(err.Error())
					}

					valueString := string(valueBytes)
					providedInputSourcer.SourceReturns(&valueString, true)

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: new(clioutputFakes.FakeCliOutput),
					}

					/* act */
					actualOutputs, err := objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					Expect(err).To(BeNil())
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
	})
})
