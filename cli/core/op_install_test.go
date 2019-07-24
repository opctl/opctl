package core

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/util/cliexiter"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("core", func() {
	Context("OpInstall", func() {
		It("should call dataResolver w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"
			providedPullCreds := &types.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			fakeDataResolver := new(fakeDataResolver)
			fakeDataResolver.ResolveReturns(new(data.FakeHandle))

			objectUnderTest := _core{
				opInstaller:  new(op.FakeInstaller),
				dataResolver: fakeDataResolver,
			}

			/* act */
			objectUnderTest.OpInstall(
				context.Background(),
				"dummyPath",
				providedPkgRef,
				providedPullCreds.Username,
				providedPullCreds.Password,
			)

			/* assert */
			actualPkgRef,
				actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)

			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(Equal(providedPullCreds))
		})
		It("should call pkg.Install w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedPath := "dummyPath"

			fakeHandle := new(data.FakeHandle)

			fakeDataResolver := new(fakeDataResolver)
			fakeDataResolver.ResolveReturns(fakeHandle)

			fakeOpInstaller := new(op.FakeInstaller)

			objectUnderTest := _core{
				opInstaller:  fakeOpInstaller,
				dataResolver: fakeDataResolver,
			}

			/* act */
			objectUnderTest.OpInstall(
				providedCtx,
				providedPath,
				"dummyPkgRef",
				"dummyUsername",
				"dummyPassword",
			)

			/* assert */
			actualContext,
				actualPath,
				actualHandle := fakeOpInstaller.InstallArgsForCall(0)

			Expect(actualContext).To(Equal(providedCtx))
			Expect(actualPath).To(Equal(providedPath))
			Expect(actualHandle).To(Equal(fakeHandle))
		})
		Context("pkg.Install errs", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeOpInstaller := new(op.FakeInstaller)

				expectedError := errors.New("dummyError")
				fakeOpInstaller.InstallReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					cliExiter:    fakeCliExiter,
					opInstaller:  fakeOpInstaller,
					dataResolver: new(fakeDataResolver),
				}

				/* act */
				objectUnderTest.OpInstall(
					context.Background(),
					"",
					"",
					"",
					"",
				)

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

			})
		})
	})
})
