package boolean

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate", func() {
	Context("value.Boolean nil", func() {
		It("should return expected errors", func() {

			/* arrange */
			providedValue := &model.Value{}

			expectedErrors := []error{
				errors.New("boolean required"),
			}

			/* act */
			actualErrors := Validate(providedValue)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))

		})
	})
	Context("value.Boolean not nil", func() {
		It("should return no errors", func() {

			/* arrange */
			providedValueBoolean := true
			providedValue := &model.Value{
				Boolean: &providedValueBoolean,
			}

			expectedErrors := []error{}

			/* act */
			actualErrors := Validate(providedValue)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))

		})
	})

})
