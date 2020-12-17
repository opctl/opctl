package core

import (
	"context"
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	dataresolver "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("Lser", func() {
	Context("Ls", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedDirRef := "dummyDirRef"

			fakeOpHandle := new(FakeDataHandle)
			fakeOpHandle.ListDescendantsReturns(nil, nil)

			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

			objectUnderTest := _lsInvoker{
				dataResolver: fakeDataResolver,
				writer:       os.Stdout,
			}

			/* act */
			err := objectUnderTest.Ls(
				context.Background(),
				providedDirRef,
			)

			/* assert */
			actualDirRef,
				actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)
			Expect(err).To(BeNil())
			Expect(actualDirRef).To(Equal(providedDirRef))
			Expect(actualPullCreds).To(BeNil())
		})
		Context("opLister.List errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.ListDescendantsReturns(nil, errors.New(""))

				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _lsInvoker{
					dataResolver: fakeDataResolver,
					writer:       os.Stdout,
				}

				/* act */
				err := objectUnderTest.Ls(
					context.Background(),
					"dummyDirRef",
				)

				/* assert */
				Expect(err).To(MatchError(""))
			})
		})
	})
})
