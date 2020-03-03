package core

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/fakes"
	"os"
)

var _ = Context("Lser", func() {
	Context("Ls", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedDirRef := "dummyDirRef"

			fakeDataResolver := new(dataresolver.Fake)

			fakeOpLister := new(FakeLister)
			// err to trigger immediate return
			fakeOpLister.ListReturns(nil, errors.New("dummyError"))

			fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

			objectUnderTest := _lsInvoker{
				dataResolver: fakeDataResolver,
				opLister:     fakeOpLister,
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
		It("should call opLister.List w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedDirRef := "dummyDirRef"

			fakeDataResolver := new(dataresolver.Fake)
			fakeDataHandle := new(FakeDataHandle)
			fakeDataResolver.ResolveReturns(fakeDataHandle)

			fakeOpLister := new(FakeLister)
			// err to trigger immediate return
			fakeOpLister.ListReturns(nil, errors.New("dummyError"))

			fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

			objectUnderTest := _lsInvoker{
				dataResolver: fakeDataResolver,
				opLister:     fakeOpLister,
				cliExiter:    fakeCliExiter,
				writer:       os.Stdout,
			}

			/* act */
			objectUnderTest.Ls(
				providedCtx,
				providedDirRef,
			)

			/* assert */
			actualCtx,
				actualDataHandle := fakeOpLister.ListArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualDataHandle).To(Equal(fakeDataHandle))
		})
		Context("opLister.List errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeOpLister := new(FakeLister)
				expectedError := errors.New("dummyError")
				fakeOpLister.ListReturns(nil, expectedError)

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _lsInvoker{
					dataResolver: new(dataresolver.Fake),
					opLister:     fakeOpLister,
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
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
