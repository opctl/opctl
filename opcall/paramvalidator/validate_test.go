package paramvalidator

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Validate", func() {
	objectUnderTest := New()
	Context("invoked w/ nil param", func() {
		It("should return expected error", func() {
			/* arrange */
			providedValue := &model.Data{}

			expectedErrs := []error{errors.New("Validate required")}

			/* act */
			actualErrs := objectUnderTest.Validate(providedValue, nil)
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
})
