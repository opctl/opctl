package core

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("core", func() {
	Context("ResolveData", func() {
		It("should call data.Resolve w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedOpRef := "dummyOpRef"
			providedPullCreds := &types.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			fakeData := new(data.Fake)

			expectedPkgProviders := []data.Provider{
				new(data.FakeProvider),
				new(data.FakeProvider),
			}
			fakeData.NewFSProviderReturns(expectedPkgProviders[0])
			fakeData.NewGitProviderReturns(expectedPkgProviders[1])

			objectUnderTest := _core{
				data: fakeData,
			}

			/* act */
			objectUnderTest.ResolveData(
				providedCtx,
				providedOpRef,
				providedPullCreds,
			)

			/* assert */
			actualCtx,
				actualOpRef,
				actualPkgProviders := fakeData.ResolveArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualOpRef).To(Equal(providedOpRef))
			Expect(actualPkgProviders).To(ConsistOf(expectedPkgProviders))
		})
	})
})
