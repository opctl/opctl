package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("core", func() {
	Context("ListPkgContents", func() {
		It("should call pkg.ListContents w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"

			fakePkg := new(pkg.Fake)

			objectUnderTest := _core{
				pkg: fakePkg,
			}

			/* act */
			objectUnderTest.ListPkgContents(providedPkgRef)

			/* assert */
			actualPkgRef := fakePkg.ListContentsArgsForCall(0)
			Expect(actualPkgRef).To(Equal(providedPkgRef))
		})
	})
})
