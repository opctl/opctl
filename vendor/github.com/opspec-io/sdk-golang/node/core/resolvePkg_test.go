package core

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("core", func() {
	Context("ResolvePkg", func() {
		It("should call pkg.Resolve w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
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
				providedCtx,
				providedPkgRef,
				providedPullCreds,
			)

			/* assert */
			actualCtx,
				actualPkgRef,
				actualPkgProviders := fakePkg.ResolveArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPkgProviders).To(ConsistOf(expectedPkgProviders))
		})
	})
})
