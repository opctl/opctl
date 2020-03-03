package params

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/internal/fakes"
)

var _ = Context("interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {

		It("should call paramsCoercer.Coerce w/ expected args", func() {
			/* arrange */
			providedScope := map[string]*model.Value{"dummyArgName": new(model.Value)}
			providedParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedOpScratchDir := "providedOpScratchDir"

			fakeParamsCoercer := new(FakeCoercer)

			objectUnderTest := _interpreter{
				paramsCoercer:   fakeParamsCoercer,
				paramsDefaulter: new(FakeDefaulter),
				paramsValidator: new(FakeValidator),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedParams,
				"providedOpPath",
				providedOpScratchDir,
			)

			/* assert */
			actualScope,
				actualParams,
				actualOpScratchDir := fakeParamsCoercer.CoerceArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualParams).To(Equal(providedParams))
			Expect(actualOpScratchDir).To(Equal(providedOpScratchDir))

		})

		It("should call paramsDefaulter.Default w/ expected args", func() {
			/* arrange */
			providedScope := map[string]*model.Value{"dummyArgName": new(model.Value)}
			providedParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedOpPath := "dummyOpPath"

			fakeParamsCoercer := new(FakeCoercer)
			coercedScope := map[string]*model.Value{"dummyArgName": new(model.Value)}

			fakeParamsCoercer.CoerceReturns(coercedScope, nil)

			fakeParamsDefaulter := new(FakeDefaulter)

			objectUnderTest := _interpreter{
				paramsCoercer:   fakeParamsCoercer,
				paramsDefaulter: fakeParamsDefaulter,
				paramsValidator: new(FakeValidator),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedParams,
				providedOpPath,
				"dummyOpScratchDir",
			)

			/* assert */
			actualScope,
				actualParams,
				actualOpPath := fakeParamsDefaulter.DefaultArgsForCall(0)

			Expect(actualScope).To(Equal(coercedScope))
			Expect(actualParams).To(Equal(providedParams))
			Expect(actualOpPath).To(Equal(providedOpPath))

		})
		It("should call paramsValidator.Validate w/ expected args & return expected result", func() {
			/* arrange */
			providedParams := map[string]*model.Param{"dummyParamName": new(model.Param)}
			providedOpPath := "dummyOpPath"

			fakeParamsDefaulter := new(FakeDefaulter)
			defaultedScope := map[string]*model.Value{"dummyArgName": new(model.Value)}

			fakeParamsDefaulter.DefaultReturns(defaultedScope)

			fakeParamsValidator := new(FakeValidator)
			validateErr := errors.New("dummyErr")
			fakeParamsValidator.ValidateReturns(validateErr)

			objectUnderTest := _interpreter{
				paramsCoercer:   new(FakeCoercer),
				paramsDefaulter: fakeParamsDefaulter,
				paramsValidator: fakeParamsValidator,
			}

			/* act */
			actualOutputs, actualErr := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedParams,
				providedOpPath,
				"dummyOpScratchDir",
			)

			/* assert */
			actualScope,
				actualParams := fakeParamsValidator.ValidateArgsForCall(0)
			Expect(actualScope).To(Equal(defaultedScope))
			Expect(actualParams).To(Equal(providedParams))

			Expect(actualOutputs).To(Equal(defaultedScope))
			Expect(actualErr).To(Equal(validateErr))

		})
	})
})
