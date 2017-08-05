package inputs

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("Validate", func() {
	objectUnderTest := newValidator()
	Context("param.Socket not nil", func() {
		Context("value.Socket nil", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Value{}
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
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
				providedValue := &model.Value{
					Socket: &providedValueSocket,
				}
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
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
