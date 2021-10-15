package cliparamsatisfier

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/clioutput"
	. "github.com/opctl/opctl/cli/internal/cliparamsatisfier/internal/fakes"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("parameterSatisfier", func() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cliOutput := clioutput.New(clicolorer.New(), io.Discard, io.Discard)

	Context("New", func() {
		It("should return truthy result", func() {
			/* arrange/act */
			actual := New(
				cliOutput,
			)

			/* assert */
			Expect(actual).To(Not(BeNil()))
		})
	})
	Context("Satisfy", func() {
		It("should call inputSourcer.Source w/ expected args for each input", func() {
			/* arrange */
			providedInputSourcer := new(FakeInputSourcer)
			providedInputSourcer.SourceReturns(nil, true)
			providedInputs := map[string]*model.ParamSpec{
				"input1": {String: &model.StringParamSpec{}},
				"input2": {String: &model.StringParamSpec{}},
			}

			expectedInputNames := map[string]struct{}{
				"input1": {},
				"input2": {},
			}

			objectUnderTest := _CLIParamSatisfier{
				cliOutput: cliOutput,
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
					providedInputs := map[string]*model.ParamSpec{
						input1Name: {Array: &model.ArrayParamSpec{}},
					}

					expectedOutputs := map[string]*model.Value{
						input1Name: {
							Array: new([]interface{}),
						},
					}

					valueBytes, err := json.Marshal(expectedOutputs[input1Name].Array)
					if err != nil {
						Fail(err.Error())
					}

					valueString := string(valueBytes)
					providedInputSourcer.SourceReturns(&valueString, true)

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: cliOutput,
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

					providedInputs := map[string]*model.ParamSpec{
						inputIdentifier: {Boolean: &model.BooleanParamSpec{}},
					}

					valueBool := true
					valueString := strconv.FormatBool(valueBool)
					providedInputSourcer.SourceReturns(&valueString, true)

					expectedOutputs := map[string]*model.Value{
						inputIdentifier: {Boolean: &valueBool},
					}

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: cliOutput,
					}

					/* act */
					actualOutputs, err := objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					Expect(err).To(BeNil())
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param.Dir isn't nil", func() {
			Context("value isn't nil", func() {
				It("should return expected outputs", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)
					inputIdentifier := "inputIdentifier"

					providedInputs := map[string]*model.ParamSpec{
						inputIdentifier: {Dir: &model.DirParamSpec{}},
					}

					valueDir := wd
					_, err := filepath.Abs(valueDir)
					if err != nil {
						panic(err)
					}

					providedInputSourcer.SourceReturns(&valueDir, true)

					expectedOutputs := map[string]*model.Value{
						inputIdentifier: {Dir: &valueDir},
					}

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: cliOutput,
					}

					/* act */
					actualOutputs, err := objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

					/* assert */
					Expect(err).To(BeNil())
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param.File isn't nil", func() {
			Context("value isn't nil", func() {
				It("should return expected outputs", func() {
					/* arrange */
					providedInputSourcer := new(FakeInputSourcer)
					inputIdentifier := "inputIdentifier"

					providedInputs := map[string]*model.ParamSpec{
						inputIdentifier: {File: &model.FileParamSpec{}},
					}

					valueFile := filepath.Join(wd, "inputSourcer.go")
					_, err := filepath.Abs(valueFile)
					if err != nil {
						panic(err)
					}

					providedInputSourcer.SourceReturns(&valueFile, true)

					expectedOutputs := map[string]*model.Value{
						inputIdentifier: {File: &valueFile},
					}

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: cliOutput,
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

					providedInputs := map[string]*model.ParamSpec{
						inputIdentifier: {Number: &model.NumberParamSpec{}},
					}

					valueNumber := 1.1
					valueString := fmt.Sprintf("%v", valueNumber)
					providedInputSourcer.SourceReturns(&valueString, true)

					expectedOutputs := map[string]*model.Value{
						inputIdentifier: {Number: &valueNumber},
					}

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: cliOutput,
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
					providedInputs := map[string]*model.ParamSpec{
						input1Name: {Object: &model.ObjectParamSpec{}},
					}

					expectedOutputs := map[string]*model.Value{
						input1Name: {
							Object: new(map[string]interface{}),
						},
					}

					valueBytes, err := json.Marshal(expectedOutputs[input1Name].Object)
					if err != nil {
						Fail(err.Error())
					}

					valueString := string(valueBytes)
					providedInputSourcer.SourceReturns(&valueString, true)

					objectUnderTest := _CLIParamSatisfier{
						cliOutput: cliOutput,
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
