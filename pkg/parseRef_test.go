package pkg

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("_Pkg", func() {
	Describe("ParseRef", func() {
		Describe("url.Parse errors", func() {
			It("should error", func() {
				/* arrange */
				providedPkgRef := "::"

				objectUnderTest := _Pkg{}

				/* act */
				_, actualErr := objectUnderTest.ParseRef(providedPkgRef)

				/* assert */
				Expect(actualErr).To(Not(BeNil()))
			})
		})
		Describe("url.Parse doesn't error", func() {
			It("should return expected PkgRef", func() {
				/* arrange */
				providedFullyQualifiedPkgName := "somehost.com/path/pkgName"
				providedPkgVersion := "0.0.0"
				providedPkgRef := fmt.Sprintf("%v#%v", providedFullyQualifiedPkgName, providedPkgVersion)
				expectedPkgRef := &PkgRef{
					FullyQualifiedName: providedFullyQualifiedPkgName,
					Version:            providedPkgVersion,
				}
				objectUnderTest := _Pkg{}

				/* act */
				actualPkgRef, actualErr := objectUnderTest.ParseRef(providedPkgRef)

				/* assert */
				Expect(actualPkgRef).To(Equal(expectedPkgRef))
				Expect(actualErr).To(BeNil())

			})
		})
	})
})
