package core

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"os"
	"path/filepath"
)

var _ = Context("pkgLs", func() {
	Context("Execute", func() {
		Context("ios.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeIOS := new(ios.Fake)
				expectedError := errors.New("dummyError")
				fakeIOS.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       new(pkg.Fake),
					cliExiter: fakeCliExiter,
					os:        fakeIOS,
					writer:    os.Stdout,
				}

				/* act */
				objectUnderTest.PkgLs("")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("ios.Getwd doesn't error", func() {
			It("should call pkg.ListOps w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				providedPath := "dummyPath"
				wdReturnedFromIOS := "dummyWorkDir"

				fakeIOS := new(ios.Fake)
				fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)
				expectedPath := filepath.Join(wdReturnedFromIOS, providedPath)

				objectUnderTest := _core{
					pkg:    fakePkg,
					os:     fakeIOS,
					writer: os.Stdout,
				}

				/* act */
				objectUnderTest.PkgLs(providedPath)

				/* assert */

				Expect(fakePkg.ListOpsArgsForCall(0)).To(Equal(expectedPath))
			})
			Context("pkg.ListOps errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					expectedError := errors.New("dummyError")
					fakePkg.ListOpsReturns([]*model.PkgManifest{}, expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						os:        new(ios.Fake),
						writer:    os.Stdout,
					}

					/* act */
					objectUnderTest.PkgLs("dummyPath")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
				})
			})
		})
	})
})
