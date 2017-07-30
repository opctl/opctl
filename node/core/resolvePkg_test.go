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
			providedPullCreds := &model.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			fakePkg := new(pkg.Fake)

			expectedPkgProviders := []pkg.Provider{
				new(pkg.FakeProvider),
				new(pkg.FakeProvider),
			}
			fakePkg.NewFSProviderReturns(expectedPkgProviders[0])
			fakePkg.NewGitProviderReturns(expectedPkgProviders[1])

			objectUnderTest := _core{
				pkg: fakePkg,
			}

			/* act */
			objectUnderTest.ResolvePkg(
				providedPkgRef,
				providedPullCreds,
			)

			/* assert */
			actualPkgRef, actualPkgProviders := fakePkg.ResolveArgsForCall(0)
			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPkgProviders).To(ConsistOf(expectedPkgProviders))
		})
	})
})
