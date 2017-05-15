package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"path/filepath"
)

var _ = Describe("Validate", func() {
	It("should call manifestValidator.Validate w/ expected args & return result", func() {
		/* arrange */
		providedPkgPath := "dummyPkgPath"

		expectedPath := filepath.Join(providedPkgPath, OpDotYmlFileName)

		expectedErrs := []error{
			errors.New("dummyErr1"),
			errors.New("dummyErr2"),
		}
		fakeManifest := new(manifest.Fake)

		fakeManifest.ValidateReturns(expectedErrs)

		objectUnderTest := _Pkg{
			manifest: fakeManifest,
		}

		/* act */
		actualErrs := objectUnderTest.Validate(providedPkgPath)

		/* assert */
		Expect(fakeManifest.ValidateArgsForCall(0)).To(Equal(expectedPath))
		Expect(actualErrs).To(Equal(expectedErrs))
	})
})
