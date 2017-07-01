package pkg

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Context("_Pkg", func() {
	Context("PkgRef", func() {
		Context("ToPath", func() {
			It("should return expected path", func() {
				/* arrange */
				providedBasePath := "/dummy/path"
				objectUnderTest := &PkgRef{
					FullyQualifiedName: "test.com/org/pkg-name",
					Version:            "0.0.0",
				}

				expectedPath := filepath.Join(
					providedBasePath,
					filepath.FromSlash(objectUnderTest.FullyQualifiedName),
					objectUnderTest.Version,
				)

				/* act */
				actualPath := objectUnderTest.ToPath(providedBasePath)

				/* assert */
				Expect(actualPath).To(Equal(expectedPath))
			})
		})
	})
	Context("ParseRef", func() {
		Context("url.Parse errors", func() {
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
		Context("url.Parse doesn't error", func() {
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
