package outputs

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("inputs", func() {
	Context("Interpret", func() {

		It("should call defaulter.Default w/ expected args", func() {
			/* arrange */
			providedOutputArgs := map[string]*model.Value{"dummyArgName": new(model.Value)}
			providedOutputParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedPkgPath := "dummyPkgPath"

			fakeDefaulter := new(fakeDefaulter)

			objectUnderTest := _interpreter{
				defaulter: fakeDefaulter,
				validator: new(fakeValidator),
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
				actualPkgPath := fakeDefaulter.DefaultArgsForCall(0)

			Expect(actualOutputArgs).To(Equal(providedOutputArgs))
			Expect(actualOutputParams).To(Equal(providedOutputParams))
			Expect(actualPkgPath).To(Equal(providedPkgPath))

		})
		It("should call validator.Validate w/ expected args & returns expected result", func() {
			/* arrange */
			providedOutputParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedPkgPath := "dummyPkgPath"

			fakeDefaulter := new(fakeDefaulter)
			defaultedOutputArgs := map[string]*model.Value{"dummyArgName": new(model.Value)}

			fakeDefaulter.DefaultReturns(defaultedOutputArgs)

			fakeValidator := new(fakeValidator)
			validateErr := errors.New("dummyErr")
			fakeValidator.ValidateReturns(validateErr)

			objectUnderTest := _interpreter{
				defaulter: fakeDefaulter,
				validator: fakeValidator,
			}

			/* act */
			actualOutputs, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedOutputParams,
				providedPkgPath,
			)

			/* assert */
			actualOutputArgs,
				actualOutputParams := fakeValidator.ValidateArgsForCall(0)
			Expect(actualOutputArgs).To(Equal(defaultedOutputArgs))
			Expect(actualOutputParams).To(Equal(providedOutputParams))

			Expect(actualOutputs).To(Equal(defaultedOutputArgs))
			Expect(actualErr).To(Equal(validateErr))

		})
	})
})
