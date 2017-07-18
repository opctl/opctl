package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("core", func() {
	Context("ResolvePkg", func() {
		It("should call pkg.Resolve w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"
			providedOpts := &pkg.ResolveOpts{
				PullCreds: &model.PullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				},
			}

			fakePkg := new(pkg.Fake)

			objectUnderTest := _core{
				pkg: fakePkg,
			}

			/* act */
			objectUnderTest.ResolvePkg(
				providedPkgRef,
				providedOpts,
			)

			/* assert */
			actualPkgRef, actualPullCreds := fakePkg.ResolveArgsForCall(0)
			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(Equal(providedOpts))
		})
	})
})
