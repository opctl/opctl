package validate

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Param", func() {
	objectUnderTest := New()
	Context("invoked w/ nil param", func() {
		It("should return expected error", func() {
			/* arrange */
			providedValue := &model.Data{}

			expectedErrs := []error{errors.New("Param required")}

			/* act */
			actualErrs := objectUnderTest.Param(providedValue, nil)
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
})
