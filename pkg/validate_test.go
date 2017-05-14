package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate", func() {
	It("should call manifestValidator.Validate w/ expected args & return result", func() {
		/* arrange */
		providedPkgPath := "dummyPkgPath"

		expectedErrs := []error{
			errors.New("dummyErr1"),
			errors.New("dummyErr2"),
		}
		fakeManifestValidator := new(fakeManifestValidator)

		fakeManifestValidator.ValidateReturns(expectedErrs)

		objectUnderTest := _Pkg{
			manifestValidator: fakeManifestValidator,
		}

		/* act */
		actualErrs := objectUnderTest.Validate(providedPkgPath)

		/* assert */
		Expect(fakeManifestValidator.ValidateArgsForCall(0)).To(Equal(providedPkgPath))
		Expect(actualErrs).To(Equal(expectedErrs))
	})
})
