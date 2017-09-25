package cliparamsatisfier

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/inputs"
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
				cliExiter: new(cliexiter.Fake),
				cliOutput: new(clioutput.Fake),
				inputs:    new(inputs.Fake),
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
							Array: []interface{}{"dummyItem"},
						},
					}

					valueBytes, err := json.Marshal(expectedValues[input1Name].Array)
					if nil != err {
						Fail(err.Error())
					}

					valueString := string(valueBytes)
					providedInputSourcer.SourceReturns(&valueString, true)

					fakeInputs := new(inputs.Fake)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter: new(cliexiter.Fake),
						cliOutput: new(clioutput.Fake),
						inputs:    fakeInputs,
					}

					/* act */
					objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					actualValues, actualParams := fakeInputs.ValidateArgsForCall(0)
					Expect(actualValues).To(Equal(expectedValues))
					Expect(actualParams).To(Equal(providedInputs))
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
							Object: map[string]interface{}{
								"prop1Name": "prop1Value",
							},
						},
					}

					valueBytes, err := json.Marshal(expectedValues[input1Name].Object)
					if nil != err {
						Fail(err.Error())
					}

					valueString := string(valueBytes)
					providedInputSourcer.SourceReturns(&valueString, true)

					fakeInputs := new(inputs.Fake)

					objectUnderTest := _CLIParamSatisfier{
						cliExiter: new(cliexiter.Fake),
						cliOutput: new(clioutput.Fake),
						inputs:    fakeInputs,
					}

					/* act */
					objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					actualValues, actualParams := fakeInputs.ValidateArgsForCall(0)
					Expect(actualValues).To(Equal(expectedValues))
					Expect(actualParams).To(Equal(providedInputs))
				})
			})
		})
	})
})
