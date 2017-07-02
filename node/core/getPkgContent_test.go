package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("core", func() {
	Context("GetPkgContent", func() {
		It("should call pkg.GetContent w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"
			providedContentPath := "dummyContentPath"

			fakePkg := new(pkg.Fake)

			objectUnderTest := _core{
				pkg: fakePkg,
			}

			/* act */
			objectUnderTest.GetPkgContent(providedPkgRef, providedContentPath)

			/* assert */
			actualPkgRef, actualContentPath := fakePkg.GetContentArgsForCall(0)
			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualContentPath).To(Equal(providedContentPath))
		})
	})
})
