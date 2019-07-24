package param

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("Validate", func() {
	objectUnderTest := NewValidator()
	Context("param.Boolean not nil", func() {
		Context("value.Boolean nil", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &types.Value{}
				providedParam := &types.Param{
					Boolean: &types.BooleanParam{},
				}

				expectedErrors := []error{
					errors.New("boolean required"),
				}

				/* act */
				actualErrors := objectUnderTest.Validate(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("value.Boolean not nil", func() {
			It("should return no errors", func() {

				/* arrange */
				providedValueBoolean := true
				providedValue := &types.Value{
					Boolean: &providedValueBoolean,
				}
				providedParam := &types.Param{
					Boolean: &types.BooleanParam{},
				}

				expectedErrors := []error{}

				/* act */
				actualErrors := objectUnderTest.Validate(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})

})
