package git

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Context("refParser", func() {
	Context("ref", func() {
		Context("ToPath", func() {
			It("should return expected path", func() {
				/* arrange */
				providedBasePath := "/dummy/path"
				objectUnderTest := &ref{
					Name:    "test.com/org/pkg-name",
					Version: "0.0.0",
				}

				expectedPath := filepath.Join(
					providedBasePath,
					filepath.FromSlash(fmt.Sprintf("%v#%v", objectUnderTest.Name, objectUnderTest.Version)),
				)

				/* act */
				actualPath := objectUnderTest.ToPath(providedBasePath)

				/* assert */
				Expect(actualPath).To(Equal(expectedPath))
			})
		})
	})
	Context("Parse", func() {
		Context("url.Parse errors", func() {
			It("should error", func() {
				/* arrange */
				providedDataRef := "::"

				/* act */
				_, actualErr := parseRef(providedDataRef)

				/* assert */
				Expect(actualErr).To(Not(BeNil()))
			})
		})
		Context("url.Parse doesn't error", func() {
			It("should return expected Ref", func() {
				/* arrange */
				providedFullyQualifiedPkgName := "somehost.com/path/pkgName"
				providedPkgVersion := "0.0.0"
				providedDataRef := fmt.Sprintf("%v#%v/some/op/path", providedFullyQualifiedPkgName, providedPkgVersion)
				expectedDataRef := &ref{
					Name:    providedFullyQualifiedPkgName,
					Version: providedPkgVersion,
				}

				/* act */
				actualDataRef, actualErr := parseRef(providedDataRef)

				/* assert */
				Expect(actualDataRef).To(Equal(expectedDataRef))
				Expect(actualErr).To(BeNil())

			})
		})
	})
})
