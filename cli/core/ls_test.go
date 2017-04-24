package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/virtual-go/vos"
	"os"
	"path/filepath"
)

var _ = Context("pkgLs", func() {
	Context("Execute", func() {
		Context("vos.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeVOS := new(vos.Fake)
				expectedError := errors.New("dummyError")
				fakeVOS.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       new(pkg.Fake),
					cliExiter: fakeCliExiter,
					vos:       fakeVOS,
					writer:    os.Stdout,
				}

				/* act */
				objectUnderTest.PkgLs("")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call pkg.List w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				providedPath := "dummyPath"
				wdReturnedFromVOS := "dummyWorkDir"

				fakeVOS := new(vos.Fake)
				fakeVOS.GetwdReturns(wdReturnedFromVOS, nil)
				expectedPath := filepath.Join(wdReturnedFromVOS, providedPath)

				objectUnderTest := _core{
					pkg:    fakePkg,
					vos:    fakeVOS,
					writer: os.Stdout,
				}

				/* act */
				objectUnderTest.PkgLs(providedPath)

				/* assert */

				Expect(fakePkg.ListArgsForCall(0)).Should(Equal(expectedPath))
			})
			Context("pkg.List errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					expectedError := errors.New("dummyError")
					fakePkg.ListReturns([]*model.PkgManifest{}, expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						vos:       new(vos.Fake),
						writer:    os.Stdout,
					}

					/* act */
					objectUnderTest.PkgLs("dummyPath")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
				})
			})
		})
	})
})
