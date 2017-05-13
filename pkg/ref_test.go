package pkg

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parsePkgRef", func() {
	Describe("invalid pkgRef", func() {
		It("should return expected result", func() {
			/* arrange */
			providedPkgRef := "invalidPkgRef"
			expectedErr := fmt.Errorf(
				"Invalid remote pkgRef: '%v'. Valid remote pkgRef's are of the form: 'host/path#semver",
				providedPkgRef,
			)

			/* act */
			_, actualErr := parsePkgRef(providedPkgRef)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Describe("valid pkgRef", func() {
		It("should return expected PkgRef", func() {
			/* arrange */
			providedFullyQualifiedPkgName := "somehost.com/path/pkgName"
			providedPkgVersion := "0.0.0"
			providedPkgRef := fmt.Sprintf("%v#%v", providedFullyQualifiedPkgName, providedPkgVersion)
			expectedPkgRef := &PkgRef{
				FullyQualifiedName: providedFullyQualifiedPkgName,
				Version:            providedPkgVersion,
			}

			/* act */
			actualPkgRef, actualErr := parsePkgRef(providedPkgRef)

			/* assert */
			Expect(actualPkgRef).To(Equal(expectedPkgRef))
			Expect(actualErr).To(BeNil())

		})
	})
})
