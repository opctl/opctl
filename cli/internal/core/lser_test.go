package core

import (
	"context"
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	dataresolver "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("Lser", func() {
	Context("Ls", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedDirRef := "dummyDirRef"

			fakeOpHandle := new(FakeDataHandle)
			fakeOpHandle.ListDescendantsReturns(nil, errors.New(""))

			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

			fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

			objectUnderTest := _lsInvoker{
				dataResolver: fakeDataResolver,
				cliExiter:    fakeCliExiter,
				writer:       os.Stdout,
			}

			/* act */
			objectUnderTest.Ls(
				context.Background(),
				providedDirRef,
			)

			/* assert */
			actualDirRef,
				actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)

			Expect(actualDirRef).To(Equal(providedDirRef))
			Expect(actualPullCreds).To(BeNil())
		})
		Context("opLister.List errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.ListDescendantsReturns(nil, errors.New(""))

				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _lsInvoker{
					dataResolver: fakeDataResolver,
					cliExiter:    fakeCliExiter,
					writer:       os.Stdout,
				}

				/* act */
				objectUnderTest.Ls(
					context.Background(),
					"dummyDirRef",
				)

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: "", Code: 1}))
			})
		})
	})
})
