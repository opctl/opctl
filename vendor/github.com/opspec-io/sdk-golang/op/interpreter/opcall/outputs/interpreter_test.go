package outputs

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall/params"
)

var _ = Context("outputs.interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {

		It("should call paramsDefaulter.Default w/ expected args", func() {
			/* arrange */
			providedOutputArgs := map[string]*model.Value{"dummyArgName": new(model.Value)}
			providedOutputParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedPkgPath := "dummyPkgPath"

			fakeParamsDefaulter := new(params.FakeDefaulter)

			objectUnderTest := _interpreter{
				paramsDefaulter: fakeParamsDefaulter,
				paramsValidator: new(params.FakeValidator),
			}

			/* act */
			objectUnderTest.Interpret(
				providedOutputArgs,
				providedOutputParams,
				providedPkgPath,
			)

			/* assert */
			actualOutputArgs,
				actualOutputParams,
				actualPkgPath := fakeParamsDefaulter.DefaultArgsForCall(0)

			Expect(actualOutputArgs).To(Equal(providedOutputArgs))
			Expect(actualOutputParams).To(Equal(providedOutputParams))
			Expect(actualPkgPath).To(Equal(providedPkgPath))

		})
		It("should call paramsValidator.Validate w/ expected args & return expected result", func() {
			/* arrange */
			providedOutputParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedPkgPath := "dummyPkgPath"

			fakeParamsDefaulter := new(params.FakeDefaulter)
			defaultedOutputArgs := map[string]*model.Value{"dummyArgName": new(model.Value)}

			fakeParamsDefaulter.DefaultReturns(defaultedOutputArgs)

			fakeParamsValidator := new(params.FakeValidator)
			validateErr := errors.New("dummyErr")
			fakeParamsValidator.ValidateReturns(validateErr)

			objectUnderTest := _interpreter{
				paramsDefaulter: fakeParamsDefaulter,
				paramsValidator: fakeParamsValidator,
			}

			/* act */
			actualOutputs, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedOutputParams,
				providedPkgPath,
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
