package params

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("Validator", func() {
	Context("NewValidator", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewValidator()).To(Not(BeNil()))
		})
	})
	Context("Validate", func() {
		It("should call paramValidator.Validate w/ expected args", func() {
			/* arrange */

			expectedName1 := "expectedName1"
			providedValues := map[string]*types.Value{
				expectedName1: new(types.Value),
			}

			providedParams := map[string]*types.Param{
				expectedName1: new(types.Param),
			}

			fakeParamValidator := new(param.FakeValidator)

			objectUnderTest := _validator{
				paramValidator: fakeParamValidator,
			}

			/* act */
			objectUnderTest.Validate(
				providedValues,
				providedParams,
			)

			/* assert */
			actualValue1,
				actualParam1 := fakeParamValidator.ValidateArgsForCall(0)

			Expect(actualValue1).To(Equal(providedValues[expectedName1]))
			Expect(actualParam1).To(Equal(providedParams[expectedName1]))
		})
	})
})
