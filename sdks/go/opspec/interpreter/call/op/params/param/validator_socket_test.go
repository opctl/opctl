package param

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("Validate", func() {
	objectUnderTest := NewValidator()
	Context("param.Socket not nil", func() {
		Context("value.Socket nil", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &types.Value{}
				providedParam := &types.Param{
					Socket: &types.SocketParam{},
				}

				expectedErrors := []error{
					errors.New("socket required"),
				}

				/* act */
				actualErrors := objectUnderTest.Validate(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("value.Socket not nil", func() {
			It("should return no errors", func() {

				/* arrange */
				providedValueSocket := "dummySocket"
				providedValue := &types.Value{
					Socket: &providedValueSocket,
				}
				providedParam := &types.Param{
					Socket: &types.SocketParam{},
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
