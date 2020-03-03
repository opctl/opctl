package git

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/provider/git/internal"
	"path/filepath"
)

var _ = Context("refParser", func() {
	Context("Ref", func() {
		Context("ToPath", func() {
			It("should return expected path", func() {
				/* arrange */
				providedBasePath := "/dummy/path"
				objectUnderTest := &internal.Ref{
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

				objectUnderTest := _refParser{}

				/* act */
				_, actualErr := objectUnderTest.Parse(providedDataRef)

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
				expectedDataRef := &internal.Ref{
					Name:    providedFullyQualifiedPkgName,
					Version: providedPkgVersion,
				}
				objectUnderTest := _refParser{}

				/* act */
				actualDataRef, actualErr := objectUnderTest.Parse(providedDataRef)

				/* assert */
				Expect(actualDataRef).To(Equal(expectedDataRef))
				Expect(actualErr).To(BeNil())

			})
		})
	})
})
