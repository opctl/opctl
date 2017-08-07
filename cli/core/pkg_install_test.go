package core

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("core", func() {
	Context("PkgInstall", func() {
		It("should call pkgResolver w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"
			providedPullCreds := &model.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			fakePkgResolver := new(fakePkgResolver)
			fakePkgResolver.ResolveReturns(new(pkg.FakeHandle))

			objectUnderTest := _core{
				pkgResolver: fakePkgResolver,
				pkg:         new(pkg.Fake),
			}

			/* act */
			objectUnderTest.PkgInstall(
				"dummyPath",
				providedPkgRef,
				providedPullCreds.Username,
				providedPullCreds.Password,
			)

			/* assert */
			actualPkgRef,
				actualPullCreds := fakePkgResolver.ResolveArgsForCall(0)
			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(Equal(providedPullCreds))
		})
		It("should call pkg.Install w/ expected args", func() {
			/* arrange */
			providedPath := "dummyPath"

			fakeHandle := new(pkg.FakeHandle)

			fakePkgResolver := new(fakePkgResolver)
			fakePkgResolver.ResolveReturns(fakeHandle)

			fakePkg := new(pkg.Fake)

			objectUnderTest := _core{
				pkg:         fakePkg,
				pkgResolver: fakePkgResolver,
			}

			/* act */
			objectUnderTest.PkgInstall(
				providedPath,
				"dummyPkgRef",
				"dummyUsername",
				"dummyPassword",
			)

			/* assert */
			actualContext, actualPath, actualHandle := fakePkg.InstallArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualPath).To(Equal(providedPath))
			Expect(actualHandle).To(Equal(fakeHandle))
		})
		Context("pkg.Install errs", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				expectedError := errors.New("dummyError")
				fakePkg.InstallReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:         fakePkg,
					cliExiter:   fakeCliExiter,
					pkgResolver: new(fakePkgResolver),
				}

				/* act */
				objectUnderTest.PkgInstall("", "", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

			})
		})
	})
})
