package outputs

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
)

var _ = Context("outputs.interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {

		It("should call paramsCoercer.Coerce w/ expected args", func() {
			/* arrange */
			providedOutputArgs := map[string]*model.Value{"dummyArgName": new(model.Value)}
			providedOutputParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedOpScratchDir := "providedOpScratchDir"

			fakeParamsCoercer := new(params.FakeCoercer)

			objectUnderTest := _interpreter{
				paramsCoercer:   fakeParamsCoercer,
				paramsDefaulter: new(params.FakeDefaulter),
				paramsValidator: new(params.FakeValidator),
			}

			/* act */
			objectUnderTest.Interpret(
				providedOutputArgs,
				providedOutputParams,
				"providedOpPath",
				providedOpScratchDir,
			)

			/* assert */
			actualOutputArgs,
				actualOutputParams,
				actualOpScratchDir := fakeParamsCoercer.CoerceArgsForCall(0)

			Expect(actualOutputArgs).To(Equal(providedOutputArgs))
			Expect(actualOutputParams).To(Equal(providedOutputParams))
			Expect(actualOpScratchDir).To(Equal(providedOpScratchDir))

		})

		It("should call paramsDefaulter.Default w/ expected args", func() {
			/* arrange */
			providedOutputArgs := map[string]*model.Value{"dummyArgName": new(model.Value)}
			providedOutputParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedOpPath := "dummyOpPath"

			fakeParamsCoercer := new(params.FakeCoercer)
			coercedOutputArgs := map[string]*model.Value{"dummyArgName": new(model.Value)}

			fakeParamsCoercer.CoerceReturns(coercedOutputArgs, nil)

			fakeParamsDefaulter := new(params.FakeDefaulter)

			objectUnderTest := _interpreter{
				paramsCoercer:   fakeParamsCoercer,
				paramsDefaulter: fakeParamsDefaulter,
				paramsValidator: new(params.FakeValidator),
			}

			/* act */
			objectUnderTest.Interpret(
				providedOutputArgs,
				providedOutputParams,
				providedOpPath,
				"dummyOpScratchDir",
			)

			/* assert */
			actualOutputArgs,
				actualOutputParams,
				actualOpPath := fakeParamsDefaulter.DefaultArgsForCall(0)

			Expect(actualOutputArgs).To(Equal(coercedOutputArgs))
			Expect(actualOutputParams).To(Equal(providedOutputParams))
			Expect(actualOpPath).To(Equal(providedOpPath))

		})
		It("should call paramsValidator.Validate w/ expected args & return expected result", func() {
			/* arrange */
			providedOutputParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedOpPath := "dummyOpPath"

			fakeParamsDefaulter := new(params.FakeDefaulter)
			defaultedOutputArgs := map[string]*model.Value{"dummyArgName": new(model.Value)}

			fakeParamsDefaulter.DefaultReturns(defaultedOutputArgs)

			fakeParamsValidator := new(params.FakeValidator)
			validateErr := errors.New("dummyErr")
			fakeParamsValidator.ValidateReturns(validateErr)

			objectUnderTest := _interpreter{
				paramsCoercer:   new(params.FakeCoercer),
				paramsDefaulter: fakeParamsDefaulter,
				paramsValidator: fakeParamsValidator,
			}

			/* act */
			actualOutputs, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedOutputParams,
				providedOpPath,
				"dummyOpScratchDir",
			)

			/* assert */
			actualOutputArgs,
				actualOutputParams := fakeParamsValidator.ValidateArgsForCall(0)
			Expect(actualOutputArgs).To(Equal(defaultedOutputArgs))
			Expect(actualOutputParams).To(Equal(providedOutputParams))

			Expect(actualOutputs).To(Equal(defaultedOutputArgs))
			Expect(actualErr).To(Equal(validateErr))

		})
	})
})
