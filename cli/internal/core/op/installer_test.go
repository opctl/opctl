package op

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/fakes"
)

var _ = Context("Installer", func() {
	Context("Install", func() {
		It("should call dataResolver w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"
			providedPullCreds := &model.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(new(FakeDataHandle))

			objectUnderTest := _installer{
				opInstaller:  new(FakeInstaller),
				dataResolver: fakeDataResolver,
			}

			/* act */
			objectUnderTest.Install(
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
		It("should call opInstaller.Install w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedPath := "dummyPath"

			fakeHandle := new(FakeDataHandle)

			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(fakeHandle)

			fakeInstaller := new(FakeInstaller)

			objectUnderTest := _installer{
				opInstaller:  fakeInstaller,
				dataResolver: fakeDataResolver,
			}

			/* act */
			objectUnderTest.Install(
				providedCtx,
				providedPath,
				"dummyPkgRef",
				"dummyUsername",
				"dummyPassword",
			)

			/* assert */
			actualContext,
				actualPath,
				actualHandle := fakeInstaller.InstallArgsForCall(0)

			Expect(actualContext).To(Equal(providedCtx))
			Expect(actualPath).To(Equal(providedPath))
			Expect(actualHandle).To(Equal(fakeHandle))
		})
		Context("pkg.Install errs", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeInstaller := new(FakeInstaller)

				expectedError := errors.New("dummyError")
				fakeInstaller.InstallReturns(expectedError)

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _installer{
					cliExiter:    fakeCliExiter,
					opInstaller:  fakeInstaller,
					dataResolver: new(dataresolver.Fake),
				}

				/* act */
				objectUnderTest.Install(
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
