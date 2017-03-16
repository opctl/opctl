package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/managepackages"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path/filepath"
)

var _ = Context("listPackages", func() {
	Context("Execute", func() {
		Context("vos.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeVos := new(vos.Fake)
				expectedError := errors.New("dummyError")
				fakeVos.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					managePackages: new(managepackages.Fake),
					cliExiter:      fakeCliExiter,
					vos:            fakeVos,
					writer:         os.Stdout,
				}

				/* act */
				objectUnderTest.ListPackages("")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call managepackages.ListPackagesInDir w/ expected args", func() {
				/* arrange */
				fakeManagePackages := new(managepackages.Fake)

				providedPath := "dummyPath"
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.Fake)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)
				expectedPath := filepath.Join(wdReturnedFromVos, providedPath)

				objectUnderTest := _core{
					managePackages: fakeManagePackages,
					vos:            fakeVos,
					writer:         os.Stdout,
				}

				/* act */
				objectUnderTest.ListPackages(providedPath)

				/* assert */

				Expect(fakeManagePackages.ListPackagesInDirArgsForCall(0)).Should(Equal(expectedPath))
			})
			Context("managepackages.ListPackagesInDir errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeManagePackages := new(managepackages.Fake)
					expectedError := errors.New("dummyError")
					fakeManagePackages.ListPackagesInDirReturns([]*model.PackageView{}, expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						managePackages: fakeManagePackages,
						cliExiter:      fakeCliExiter,
						vos:            new(vos.Fake),
						writer:         os.Stdout,
					}

					/* act */
					objectUnderTest.ListPackages("dummyPath")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
				})
			})
		})
	})
})
