package op

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/model"
	op "github.com/opctl/opctl/sdks/go/opspec"
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
			fakeDataResolver.ResolveReturns(new(data.FakeHandle))

			objectUnderTest := _installer{
				opInstaller:  new(op.FakeInstaller),
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

			fakeHandle := new(data.FakeHandle)

			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(fakeHandle)

			fakeInstaller := new(op.FakeInstaller)

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
				fakeInstaller := new(op.FakeInstaller)

				expectedError := errors.New("dummyError")
				fakeInstaller.InstallReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

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
