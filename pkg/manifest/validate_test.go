package manifest

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Validate", func() {
	It("should call validator.Validate w/ expected args & return result", func() {
		/* arrange */
		providedPkgPath := "dummyPkgPath"

		expectedErrs := []error{
			errors.New("dummyErr1"),
			errors.New("dummyErr2"),
		}
		fakeValidator := new(fakeValidator)

		fakeValidator.ValidateReturns(expectedErrs)

		objectUnderTest := _Manifest{
			validator: fakeValidator,
		}

		/* act */
		actualErrs := objectUnderTest.Validate(providedPkgPath)

		/* assert */
		Expect(fakeValidator.ValidateArgsForCall(0)).To(Equal(providedPkgPath))
		Expect(actualErrs).To(Equal(expectedErrs))
	})
})
