package socket

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate", func() {

	Context("value.Socket nil", func() {
		It("should return expected errors", func() {

			/* arrange */
			providedValue := &model.Value{}

			expectedErrors := []error{
				errors.New("socket required"),
			}

			/* act */
			actualErrors := Validate(providedValue)

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

			expectedErrors := []error{}

			/* act */
			actualErrors := Validate(providedValue)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))

		})
	})
})
