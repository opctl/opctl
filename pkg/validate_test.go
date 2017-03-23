package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate", func() {
	It("should call validator.Validate w/ expected args & return result", func() {
		/* arrange */
		providedPkgRef := "dummyPkgRef"

		expectedErrs := []error{
			errors.New("dummyErr1"),
			errors.New("dummyErr2"),
		}
		fakeValidator := new(fakeValidator)

		fakeValidator.ValidateReturns(expectedErrs)

		objectUnderTest := pkg{
			validator: fakeValidator,
		}

		/* act */
		actualErrs := objectUnderTest.Validate(providedPkgRef)

		/* assert */
		Expect(fakeValidator.ValidateArgsForCall(0)).To(Equal(providedPkgRef))
		Expect(actualErrs).To(Equal(expectedErrs))
	})
})
