package param

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate", func() {
	objectUnderTest := NewValidator()
	Context("invoked w/ nil param", func() {
		It("should return expected error", func() {
			/* arrange */
			expectedErrs := []error{errors.New("required")}

			/* act */
			actualErrs := objectUnderTest.Validate(
				nil,
				&model.Param{},
			)
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
})
